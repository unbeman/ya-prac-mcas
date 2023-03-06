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

type RepositoryConfig struct {
	FileStorage *FileStorageConfig
}

type FileStorageConfig struct {
	Interval time.Duration `env:"STORE_INTERVAL"`
	File     string        `env:"STORE_FILE"`
	Restore  bool          `env:"RESTORE"`
}

func (cfg *FileStorageConfig) String() string {
	return fmt.Sprintf("[Interval: %v; File: %v; Restore: %v]", cfg.Interval, cfg.File, cfg.Restore)
}

func newFileStorageConfig() *FileStorageConfig {
	return &FileStorageConfig{
		Interval: StoreIntervalDefault,
		File:     StoreFileDefault,
		Restore:  RestoreDefault,
	}
}

type ServerConfig struct {
	Address    string `env:"ADDRESS"`
	Repository RepositoryConfig
	Logger     LoggerConfig
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
		cfg.Repository.FileStorage.Restore = *restore
		cfg.Repository.FileStorage.Interval = *storeInterval
		cfg.Repository.FileStorage.File = *storeFile
		cfg.Logger.Level = *logLevel
	}
}

func NewServerConfig(options ...ServerOption) *ServerConfig {
	cfg := &ServerConfig{
		Address:    AddressDefault,
		Repository: RepositoryConfig{FileStorage: newFileStorageConfig()},
		Logger:     newLoggerConfig(),
	}
	for _, option := range options {
		option(cfg)
	}

	if len(cfg.Repository.FileStorage.File) == 0 {
		cfg.Repository.FileStorage = nil
	}
	return cfg
}
