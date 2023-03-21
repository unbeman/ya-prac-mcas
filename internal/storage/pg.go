package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type Statements struct {
	AddCounter *sql.Stmt
	GetCounter *sql.Stmt
	SetGauge   *sql.Stmt
	GetGauge   *sql.Stmt
}

type postgresRepository struct {
	connection *sql.DB
	statements Statements
}

func (p *postgresRepository) AddCounter(ctx context.Context, name string, delta int64) (metrics.Counter, error) {
	row := p.statements.AddCounter.QueryRowContext(ctx, name, delta)
	err := row.Scan(&delta)
	if err != nil {
		return nil, err
	}
	counter := metrics.NewCounter(name)
	counter.Add(delta)
	return counter, nil
}

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
	counter := metrics.NewCounter(name)
	counter.Add(value)
	return counter, nil
}

func (p *postgresRepository) SetGauge(ctx context.Context, name string, value float64) (metrics.Gauge, error) {
	_, err := p.statements.SetGauge.ExecContext(ctx, name, value)
	if err != nil {
		log.Error(err)
	}
	gauge := metrics.NewGauge(name)
	gauge.Set(value)
	return gauge, nil
}

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
	gauge := metrics.NewGauge(name)
	gauge.Set(value)
	return gauge, nil
}

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
		gauge := metrics.NewGauge(name)
		gauge.Set(value)
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
		counter := metrics.NewCounter(name)
		counter.Add(value)
		metricSlice = append(metricSlice, counter)
	}

	err = rowsCounter.Err()
	if err != nil {
		return nil, err
	}

	return metricSlice, nil
}

func (p *postgresRepository) AddCounters(ctx context.Context, slice []metrics.Counter) error {
	transaction, err := p.connection.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()

	stmt := transaction.StmtContext(ctx, p.statements.AddCounter)
	for _, counter := range slice {
		_, err := stmt.ExecContext(ctx, counter.GetName(), counter.Value())
		if err != nil {
			return err
		}
	}
	return transaction.Commit()
}

func (p *postgresRepository) SetGauges(ctx context.Context, slice []metrics.Gauge) error {
	transaction, err := p.connection.Begin()
	if err != nil {
		return err
	}
	defer transaction.Rollback()
	stmt := transaction.StmtContext(ctx, p.statements.SetGauge)
	for _, gauge := range slice {
		_, err := stmt.ExecContext(ctx, gauge.GetName(), gauge.Value())
		if err != nil {
			return err
		}
	}
	return transaction.Commit()
}

func (p *postgresRepository) Ping() error {
	return p.connection.Ping()
}

func (p *postgresRepository) Shutdown() error {
	err := p.connection.Close()
	log.Infoln("db conn closed")
	return err
}

func (p *postgresRepository) createSchemaIfNotExist(filename string) {
	text, err := os.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	script := string(text)
	_, err = p.connection.Exec(script)
	if err != nil {
		log.Fatalln(err)
	}
}

func NewPostgresRepository(cfg configs.PostgresConfig) (*postgresRepository, error) {
	connection, err := sql.Open("pgx", cfg.DSN) //TODO: настроить пул коннектов, таймауты
	if err != nil {
		return nil, err
	}
	pg := &postgresRepository{connection: connection, statements: NewStatements(connection)}
	pg.createSchemaIfNotExist(cfg.SchemaFile)

	return pg, nil
}

func NewStatements(conn *sql.DB) Statements {
	s := Statements{}
	s.GetCounter = newStatement(conn, "SELECT value FROM counter WHERE name=$1")
	s.AddCounter = newStatement(conn, "INSERT into counter values ($1, $2) ON CONFLICT (name) DO UPDATE set value=counter.value+$2 where counter.name=$1 RETURNING value")
	s.GetGauge = newStatement(conn, "SELECT value FROM gauge WHERE name=$1")
	s.SetGauge = newStatement(conn, "INSERT into gauge values ($1, $2) ON CONFLICT (name) DO UPDATE set value=$2 where gauge.name=$1")
	return s
}

func newStatement(conn *sql.DB, query string) *sql.Stmt {
	stmt, err := conn.Prepare(query)
	if err != nil {
		log.Fatal(err)
	}
	return stmt
}