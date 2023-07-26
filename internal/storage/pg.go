package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

// Statements contains all necessary statements.
type Statements struct {
	AddCounter *sql.Stmt
	GetCounter *sql.Stmt
	SetGauge   *sql.Stmt
	GetGauge   *sql.Stmt
}

// NewStatements creates Statements.
func NewStatements(conn *sql.DB) (Statements, error) {
	var err error
	s := Statements{}
	s.GetCounter, err = conn.Prepare("SELECT value FROM counter WHERE name=$1")
	if err != nil {
		return s, err
	}
	s.AddCounter, err = conn.Prepare("INSERT into counter values ($1, $2) ON CONFLICT (name) DO UPDATE set value=counter.value+$2 where counter.name=$1 RETURNING value")
	if err != nil {
		return s, err
	}
	s.GetGauge, err = conn.Prepare("SELECT value FROM gauge WHERE name=$1")
	if err != nil {
		return s, err
	}
	s.SetGauge, err = conn.Prepare("INSERT into gauge values ($1, $2) ON CONFLICT (name) DO UPDATE set value=$2 where gauge.name=$1")
	if err != nil {
		return s, err
	}
	return s, nil
}

// postgresRepository implements Repository interface and describes PostgresSQL connection and prepared queries.
type postgresRepository struct {
	connection *sql.DB
	statements Statements
}

// NewPostgresRepository creates and configured postgresRepository,
// including migrations and statements preparation.
func NewPostgresRepository(cfg configs.PostgresConfig) (*postgresRepository, error) {
	connection, err := sql.Open("pgx", cfg.DSN) //TODO: настроить пул коннектов, таймауты
	if err != nil {
		return nil, err
	}
	pg := &postgresRepository{connection: connection}
	err = pg.migrate(cfg.MigrationDir)
	if err != nil {
		return nil, err
	}
	pg.statements, err = NewStatements(connection)
	if err != nil {
		return nil, err
	}
	return pg, nil
}

// AddCounter increases by delta counter and return metrics.Counter,
func (p *postgresRepository) AddCounter(ctx context.Context, name string, delta int64) (metrics.Counter, error) {
	row := p.statements.AddCounter.QueryRowContext(ctx, name, delta)
	err := row.Scan(&delta)
	if err != nil {
		return nil, err
	}
	counter := metrics.NewCounter(name, delta)
	return counter, nil
}

// GetCounter return counter by name.
func (p *postgresRepository) GetCounter(ctx context.Context, name string) (metrics.Counter, error) {
	row := p.statements.GetCounter.QueryRowContext(ctx, name)
	var value int64
	err := row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("counter (%v) %w", name, ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	counter := metrics.NewCounter(name, value)
	return counter, nil
}

// SetGauge set gauge metric to value and return metrics.Gauge.
func (p *postgresRepository) SetGauge(ctx context.Context, name string, value float64) (metrics.Gauge, error) {
	_, err := p.statements.SetGauge.ExecContext(ctx, name, value)
	if err != nil {
		log.Error(err)
	}
	gauge := metrics.NewGauge(name, value)
	return gauge, nil
}

// GetGauge return metrics.Gauge by name.
func (p *postgresRepository) GetGauge(ctx context.Context, name string) (metrics.Gauge, error) {
	row := p.statements.GetGauge.QueryRowContext(ctx, name)
	var value float64
	err := row.Scan(&value)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("gauge (%v) %w", name, ErrNotFound)
	}
	if err != nil {
		return nil, err
	}
	gauge := metrics.NewGauge(name, value)
	return gauge, nil
}

// GetAll return slice of all saved metrics.
func (p *postgresRepository) GetAll(ctx context.Context) ([]metrics.Metric, error) {
	metricSlice := make([]metrics.Metric, 0)

	queryGauge := "SELECT name, value FROM gauge"
	queryCounter := "SELECT name, value FROM counter"

	rowsGauge, err := p.connection.QueryContext(ctx, queryGauge)
	if err != nil {
		return nil, err
	}
	defer rowsGauge.Close()

	rowsCounter, err := p.connection.QueryContext(ctx, queryCounter)
	if err != nil {
		return nil, err
	}

	defer rowsCounter.Close()

	for rowsGauge.Next() {
		var (
			name  string
			value float64
		)
		err = rowsGauge.Scan(&name, &value)
		if err != nil {
			log.Error(err)
		}
		gauge := metrics.NewGauge(name, value)
		metricSlice = append(metricSlice, gauge)
	}
	err = rowsGauge.Err()
	if err != nil {
		return nil, err
	}

	for rowsCounter.Next() {
		var (
			name  string
			value int64
		)
		err = rowsCounter.Scan(&name, &value)
		if err != nil {
			return nil, err
		}
		counter := metrics.NewCounter(name, value)
		metricSlice = append(metricSlice, counter)
	}

	err = rowsCounter.Err()
	if err != nil {
		return nil, err
	}

	return metricSlice, nil
}

// AddCounters increase each metrics.Counter on value in slice and return slice of result.
func (p *postgresRepository) AddCounters(ctx context.Context, slice []metrics.Counter) ([]metrics.Counter, error) {
	transaction, err := p.connection.Begin()
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()

	stmt := transaction.StmtContext(ctx, p.statements.AddCounter)
	for idx, counter := range slice {
		result := stmt.QueryRowContext(ctx, counter.GetName(), counter.Value())
		if result.Err() != nil {
			return nil, err
		}

		var updatedValue int64
		err = result.Scan(&updatedValue)
		if err != nil {
			return nil, err
		}

		slice[idx].Set(updatedValue)
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}
	return slice, nil
}

// SetGauges set new value for each metrics.Gauge in slice and return the result slice.
func (p *postgresRepository) SetGauges(ctx context.Context, slice []metrics.Gauge) ([]metrics.Gauge, error) {
	transaction, err := p.connection.Begin()
	if err != nil {
		return nil, err
	}
	defer transaction.Rollback()
	stmt := transaction.StmtContext(ctx, p.statements.SetGauge)
	for _, gauge := range slice {
		_, err = stmt.ExecContext(ctx, gauge.GetName(), gauge.Value())
		if err != nil {
			return nil, err
		}
	}
	err = transaction.Commit()
	if err != nil {
		return nil, err
	}
	return slice, nil
}

// Ping checks the PG connection is alive.
func (p *postgresRepository) Ping() error {
	return p.connection.Ping()
}

// Shutdown closes the PG connection.
func (p *postgresRepository) Shutdown() error {
	err := p.connection.Close()
	log.Infoln("db conn closed")
	return err
}

// migrate makes migrate up.
func (p *postgresRepository) migrate(directory string) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}
	return goose.Up(p.connection, directory)
}
