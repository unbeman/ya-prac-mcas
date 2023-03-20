package storage

import (
	"database/sql"
	"os"

	log "github.com/sirupsen/logrus"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

type postgresRepository struct {
	connection *sql.DB
}

func (p *postgresRepository) AddCounter(name string, delta int64) metrics.Counter {
	query := "INSERT into counter values ($1, $2) ON CONFLICT (name) DO UPDATE set value=$2 where counter.name=$1 RETURNING value"
	row := p.connection.QueryRow(query, name, delta)
	err := row.Scan(&delta)
	if err != nil {
		log.Error(err)
	}
	counter := metrics.NewCounter(name)
	counter.Add(delta)
	return counter
}

func (p *postgresRepository) GetCounter(name string) metrics.Counter {
	query := "SELECT value FROM counter WHERE name=$1"
	row := p.connection.QueryRow(query, name)
	var value int64
	err := row.Scan(&value)
	if err != nil {
		log.Error(err)
		return nil
	}
	counter := metrics.NewCounter(name)
	counter.Add(value)
	return counter
}

func (p *postgresRepository) SetGauge(name string, value float64) metrics.Gauge {
	query := "INSERT into gauge values ($1, $2) ON CONFLICT (name) DO UPDATE set value=$2 where gauge.name=$1"
	_, err := p.connection.Exec(query, name, value)
	if err != nil {
		log.Error(err)
	}
	gauge := metrics.NewGauge(name)
	gauge.Set(value)
	return gauge
}

func (p *postgresRepository) GetGauge(name string) metrics.Gauge {
	query := "SELECT value FROM gauge WHERE name=$1"
	row := p.connection.QueryRow(query, name)
	var value float64
	err := row.Scan(&value)
	if err != nil {
		log.Error(err)
		return nil
	}
	gauge := metrics.NewGauge(name)
	gauge.Set(value)
	return gauge
}

func (p *postgresRepository) GetAll() []metrics.Metric {
	metricSlice := make([]metrics.Metric, 0)

	queryGauge := "SELECT name, value FROM gauge"
	queryCounter := "SELECT name, value FROM counter"

	rowsGauge, err := p.connection.Query(queryGauge)
	if err != nil {
		log.Error(err)
	}

	rowsCounter, err := p.connection.Query(queryCounter)
	if err != nil {
		log.Error(err)
	}

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

	for rowsCounter.Next() {
		var (
			name  string
			value int64
		)
		err = rowsCounter.Scan(&name, &value)
		if err != nil {
			log.Error(err)
		}
		counter := metrics.NewCounter(name)
		counter.Add(value)
		metricSlice = append(metricSlice, counter)
	}

	return metricSlice
}

func (p *postgresRepository) Ping() error {
	return p.connection.Ping()
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

func (p *postgresRepository) Shutdown() error {
	err := p.connection.Close()
	log.Infoln("db conn closed")
	return err
}

func NewPostgresRepository(cfg configs.PostgresConfig) (*postgresRepository, error) {
	sql.Drivers()
	connection, err := sql.Open("pgx", cfg.DSN)
	if err != nil {
		return nil, err
	}
	pg := &postgresRepository{connection: connection}
	pg.createSchemaIfNotExist(cfg.SchemaFile)

	return pg, nil
}
