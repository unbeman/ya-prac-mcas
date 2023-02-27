package configs

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type FileStorageConfig struct {
	Interval time.Duration `env:"STORE_INTERVAL"`
	File     string        `env:"STORE_FILE"`
	Restore  bool          `env:"RESTORE"`
}

func newFileStorageConfig() *FileStorageConfig {
	return &FileStorageConfig{
		Interval: 300 * time.Second,
		File:     "/tmp/devops-metrics-db.json",
		Restore:  true,
	}
}

type ServerConfig struct {
	Address string `env:"ADDRESS"`
	File    FileStorageConfig
}

func (cfg *ServerConfig) FromEnv() *ServerConfig {
	if err := env.Parse(cfg); err != nil {
		log.Fatalln("ServerConfig.FromEnv: can't parse env vars, reason: ", err)
	}
	return cfg
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{
		Address: "127.0.0.1:8080",
		File:    *newFileStorageConfig(),
	}
	return cfg
}
