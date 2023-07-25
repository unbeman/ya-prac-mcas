// Package configs describes applications settings.
package configs

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

// Default config settings
const (
	ProfileAddressDefault       = "127.0.0.1:8888"
	BackupIntervalDefault       = 300 * time.Second
	BackupFileDefault           = "/tmp/devops-metrics-db.json"
	RestoreDefault              = true
	DSNDefault                  = ""
	PGMigrationDirDefault       = "migrations"
	PrivateCryptoKeyPathDefault = "private.pem"
	JSONServerConfigPathDefault = ""
)

type ServerOption func(config *ServerConfig)

type PostgresConfig struct {
	DSN          string `env:"DATABASE_DSN" json:"database_dsn,omitempty"`
	MigrationDir string `env:"MIGRATION_DIR" json:"migration_dir,omitempty"`
}

func (cfg *PostgresConfig) String() string {
	return fmt.Sprintf("[DSN: %v, MigrationDir: %v]", cfg.DSN, cfg.MigrationDir)
}

func newPostgresConfig() *PostgresConfig {
	return &PostgresConfig{DSN: DSNDefault, MigrationDir: PGMigrationDirDefault}
}

type RepositoryConfig struct {
	RAMWithBackup *BackupConfig
	PG            *PostgresConfig
}

type BackupConfig struct {
	Interval time.Duration `env:"STORE_INTERVAL"`
	Restore  bool          `env:"RESTORE" json:"restore,omitempty"`
	File     string        `env:"STORE_FILE" json:"file,omitempty"`
}

func (cfg *BackupConfig) String() string {
	return fmt.Sprintf("[Interval: %v; File: %v; Restore: %v;]", cfg.Interval, cfg.File, cfg.Restore)
}

func (cfg *BackupConfig) UnmarshalJSON(data []byte) error {
	type RealCfg BackupConfig
	jCfg := struct {
		Interval string `json:"store_interval,omitempty"`
		*RealCfg
	}{
		RealCfg: (*RealCfg)(cfg),
	}

	err := json.Unmarshal(data, &jCfg)
	if err != nil {
		return err
	}
	if jCfg.Interval != "" {
		cfg.Interval, err = time.ParseDuration(jCfg.Interval)
		if err != nil {
			return err
		}
	}

	return nil
}

func newBackupConfig() *BackupConfig {
	return &BackupConfig{
		Interval: BackupIntervalDefault,
		File:     BackupFileDefault,
		Restore:  RestoreDefault,
	}
}

type ServerConfig struct {
	CollectorAddress     string `env:"ADDRESS" json:"address,omitempty"`
	HashKey              string `env:"KEY" json:"key,omitempty"`
	PrivateCryptoKeyPath string `env:"CRYPTO_KEY" json:"crypto_key,omitempty"`
	Logger               LoggerConfig
	Repository           RepositoryConfig
	ProfileAddress       string `json:"profile_address,omitempty"`
	jsonConfigPath       string
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
		privateCryptoKeyPath := flag.String("crypto-key", PrivateCryptoKeyPathDefault, "path to private key file")
		restore := flag.Bool("r", RestoreDefault, "restore metrics to file")
		storeInterval := flag.Duration("i", BackupIntervalDefault, "store interval")
		storeFile := flag.String("f", BackupFileDefault, "json file path to store metrics")
		logLevel := flag.String("l", LogLevelDefault, "log level, allowed [info, debug]")
		dsn := flag.String("d", DSNDefault, "Postgres data source name")
		jsonConfigPath := flag.String("c", cfg.jsonConfigPath, "path to json config")
		flag.Parse()
		cfg.CollectorAddress = *address
		cfg.HashKey = *key
		cfg.PrivateCryptoKeyPath = *privateCryptoKeyPath
		cfg.Repository.RAMWithBackup.Restore = *restore
		cfg.Repository.RAMWithBackup.Interval = *storeInterval
		cfg.Repository.RAMWithBackup.File = *storeFile
		cfg.Repository.PG.DSN = *dsn
		cfg.Logger.Level = *logLevel
		cfg.jsonConfigPath = *jsonConfigPath
	}
}

func FromJSON() ServerOption {
	return func(cfg *ServerConfig) {
		envPath, isSet := os.LookupEnv("CONFIG")
		if isSet {
			cfg.jsonConfigPath = envPath
		}
		if cfg.jsonConfigPath == "" {
			return
		}
		data, err := os.ReadFile(cfg.jsonConfigPath)
		if err != nil {
			log.Fatalf("can't read %v, reason: %v", cfg.jsonConfigPath, err)
		}

		err = json.Unmarshal(data, &cfg)
		if err != nil {
			log.Fatalf("can't unmarshal json config, reason: %v", err)
		}

		err = json.Unmarshal(data, &cfg.Logger)
		if err != nil {
			log.Fatalf("can't unmarshal json config, reason: %v", err)
		}

		err = json.Unmarshal(data, &cfg.Repository.PG)
		if err != nil {
			log.Fatalf("can't unmarshal json config, reason: %v", err)
		}

		err = json.Unmarshal(data, &cfg.Repository.RAMWithBackup)
		if err != nil {
			log.Fatalf("can't unmarshal json config, reason: %v", err)
		}

	}
}

func NewServerConfig(options ...ServerOption) *ServerConfig {
	cfg := &ServerConfig{
		jsonConfigPath:       JSONServerConfigPathDefault,
		CollectorAddress:     ServerAddressDefault,
		ProfileAddress:       ProfileAddressDefault,
		HashKey:              KeyDefault,
		PrivateCryptoKeyPath: PrivateCryptoKeyPathDefault,
		Logger:               newLoggerConfig(),
		Repository:           RepositoryConfig{RAMWithBackup: newBackupConfig(), PG: newPostgresConfig()},
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
