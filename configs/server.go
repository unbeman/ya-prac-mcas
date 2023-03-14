package configs

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

const (
	AddressDefault        = "127.0.0.1:8080"
	BackupIntervalDefault = 300 * time.Second
	BackupFileDefault     = "/tmp/devops-metrics-db.json"
	RestoreDefault        = true
)

type ServerOption func(config *ServerConfig)

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
	Address string `env:"ADDRESS"`
	Logger  LoggerConfig
	Backup  *BackupConfig
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
		address := flag.String("a", AddressDefault, "server address")
		restore := flag.Bool("r", RestoreDefault, "restore metrics to file")
		storeInterval := flag.Duration("i", BackupIntervalDefault, "store interval")
		storeFile := flag.String("f", BackupFileDefault, "json file path to store metrics")
		logLevel := flag.String("l", LogLevelDefault, "log level, allowed [info, debug]")
		flag.Parse()
		cfg.Address = *address
		cfg.Backup.Restore = *restore
		cfg.Backup.Interval = *storeInterval
		cfg.Backup.File = *storeFile
		cfg.Logger.Level = *logLevel
	}
}

func NewServerConfig(options ...ServerOption) *ServerConfig {
	cfg := &ServerConfig{
		Address: AddressDefault,
		Backup:  newBackupConfig(),
		Logger:  newLoggerConfig(),
	}
	for _, option := range options {
		option(cfg)
	}
	return cfg
}
