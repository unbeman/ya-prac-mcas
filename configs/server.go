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
		flag.Func("c", "path to json config", cfg.fromJSON)
		flag.Func("config", "path to json config", cfg.fromJSON)

		flag.StringVar(&cfg.CollectorAddress, "a", cfg.CollectorAddress, "server address")
		flag.StringVar(&cfg.HashKey, "k", cfg.HashKey, "key for calculating the metric hash")
		flag.StringVar(&cfg.PrivateCryptoKeyPath, "crypto-key", cfg.PrivateCryptoKeyPath, "path to private key file")
		flag.BoolVar(&cfg.Repository.RAMWithBackup.Restore, "r", cfg.Repository.RAMWithBackup.Restore, "restore metrics to file")
		flag.DurationVar(&cfg.Repository.RAMWithBackup.Interval, "i", cfg.Repository.RAMWithBackup.Interval, "store interval")
		flag.StringVar(&cfg.Repository.RAMWithBackup.File, "f", cfg.Repository.RAMWithBackup.File, "json file path to store metrics")
		flag.StringVar(&cfg.Logger.Level, "l", cfg.Logger.Level, "log level, allowed [info, debug]")
		flag.StringVar(&cfg.Repository.PG.DSN, "d", cfg.Repository.PG.DSN, "Postgres data source name")

		flag.Parse()
	}
}

func (cfg *ServerConfig) fromJSON(path string) error {
	envPath, isSet := os.LookupEnv("CONFIG")
	if !isSet {
		envPath = path
	}

	if envPath == "" {
		return nil
	}

	data, err := os.ReadFile(envPath)
	if err != nil {
		log.Fatalf("can't read %v, reason: %v", envPath, err)
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
	return nil
}

func NewServerConfig(options ...ServerOption) *ServerConfig {
	cfg := &ServerConfig{
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
