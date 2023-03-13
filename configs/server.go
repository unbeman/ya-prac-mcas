package configs

import (
	"flag"
	"fmt"
	"time"

	"github.com/caarlos0/env/v6"
	log "github.com/sirupsen/logrus"
)

const (
	AddressDefault       = "127.0.0.1:8080"
	StoreIntervalDefault = 300 * time.Second
	StoreFileDefault     = "/tmp/devops-metrics-db.json"
	RestoreDefault       = true
)

type ServerOption func(config *ServerConfig)

type FileStorageConfig struct {
	Interval time.Duration `env:"STORE_INTERVAL"`
	File     string        `env:"STORE_FILE"`
}

func (cfg *FileStorageConfig) String() string {
	return fmt.Sprintf("[Interval: %v; File: %v;]", cfg.Interval, cfg.File)
}

func newFileStorageConfig() *FileStorageConfig {
	return &FileStorageConfig{
		Interval: StoreIntervalDefault,
		File:     StoreFileDefault,
	}
}

type ServerConfig struct {
	Address     string `env:"ADDRESS"`
	Restore     bool   `env:"RESTORE"`
	FileStorage *FileStorageConfig
	Logger      LoggerConfig
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
		storeInterval := flag.Duration("i", StoreIntervalDefault, "store interval")
		storeFile := flag.String("f", StoreFileDefault, "json file path to store metrics")
		logLevel := flag.String("l", LogLevelDefault, "log level, allowed [info, debug]")
		flag.Parse()
		cfg.Address = *address
		cfg.Restore = *restore
		cfg.FileStorage.Interval = *storeInterval
		cfg.FileStorage.File = *storeFile
		cfg.Logger.Level = *logLevel
	}
}

func NewServerConfig(options ...ServerOption) *ServerConfig {
	cfg := &ServerConfig{
		Address:     AddressDefault,
		Restore:     RestoreDefault,
		FileStorage: newFileStorageConfig(),
		Logger:      newLoggerConfig(),
	}
	for _, option := range options {
		option(cfg)
	}
	return cfg
}
