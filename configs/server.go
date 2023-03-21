package configs

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

const (
	BackupIntervalDefault = 300 * time.Second
	BackupFileDefault     = "/tmp/devops-metrics-db.json"
	RestoreDefault        = true
	DSNDefault            = ""
	PGSchemaFileDefault   = "./pg-schema.sql"
)

type ServerOption func(config *ServerConfig)

type PostgresConfig struct {
	DSN        string `env:"DATABASE_DSN"`
	SchemaFile string
}

func (cfg *PostgresConfig) String() string {
	return fmt.Sprintf("[DSN: %v, schema file: %v]", cfg.DSN, cfg.SchemaFile)
}

func newPostgresConfig() *PostgresConfig {
	return &PostgresConfig{DSN: DSNDefault, SchemaFile: PGSchemaFileDefault}
}

type RepositoryConfig struct {
	RAMWithBackup *BackupConfig
	PG            *PostgresConfig
}

type BackupConfig struct {
	Interval time.Duration `env:"STORE_INTERVAL"`
	Restore  bool          `env:"RESTORE"`
	File     string        `env:"STORE_FILE"`
}

func (cfg *BackupConfig) String() string {
	return fmt.Sprintf("[Interval: %v; File: %v; Restore: %v;]", cfg.Interval, cfg.File, cfg.Restore)
}

func newBackupConfig() *BackupConfig {
	return &BackupConfig{
		Interval: BackupIntervalDefault,
		File:     BackupFileDefault,
		Restore:  RestoreDefault,
	}
}

type ServerConfig struct {
	Address    string `env:"ADDRESS"`
	Key        string `env:"KEY"`
	Logger     LoggerConfig
	Repository RepositoryConfig
}

func FromEnv() ServerOption {
	return func(cfg *ServerConfig) {
		if err := env.Parse(cfg); err != nil {
			log.Fatalln("ServerConfig.FromEnv: can't parse env vars, reason: ", err)
		}
	}
}

func FromFlags() ServerOption {
	return func(cfg *ServerConfig) {
		address := flag.String("a", ServerAddressDefault, "server address")
		key := flag.String("k", KeyDefault, "key for calculating the metric hash")
		restore := flag.Bool("r", RestoreDefault, "restore metrics to file")
		storeInterval := flag.Duration("i", BackupIntervalDefault, "store interval")
		storeFile := flag.String("f", BackupFileDefault, "json file path to store metrics")
		logLevel := flag.String("l", LogLevelDefault, "log level, allowed [info, debug]")
		dsn := flag.String("d", DSNDefault, "Postgres data source name")
		flag.Parse()
		cfg.Address = *address
		cfg.Key = *key
		cfg.Repository.RAMWithBackup.Restore = *restore
		cfg.Repository.RAMWithBackup.Interval = *storeInterval
		cfg.Repository.RAMWithBackup.File = *storeFile
		cfg.Repository.PG.DSN = *dsn
		cfg.Logger.Level = *logLevel
	}
}

func NewServerConfig(options ...ServerOption) *ServerConfig {
	cfg := &ServerConfig{
		Address:    ServerAddressDefault,
		Key:        KeyDefault,
		Logger:     newLoggerConfig(),
		Repository: RepositoryConfig{RAMWithBackup: newBackupConfig(), PG: newPostgresConfig()},
	}
	for _, option := range options {
		option(cfg)
	}

	if cfg.Repository.PG.DSN == DSNDefault { //TODO: wrap
		cfg.Repository.PG = nil
	}

	if len(cfg.Repository.RAMWithBackup.File) == 0 { //TODO: wrap
		cfg.Repository.RAMWithBackup = nil
	}

	return cfg
}
