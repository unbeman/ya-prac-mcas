package configs

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type FileHandlerConfig struct {
	Interval time.Duration `env:"STORE_INTERVAL"`
	File     string        `env:"STORE_FILE"`
	Restore  bool          `env:"RESTORE"`
}

func newFileStorageConfig() *FileHandlerConfig {
	return &FileHandlerConfig{
		Interval: 300 * time.Second,
		File:     "/tmp/devops-metrics-db.json",
		Restore:  true,
	}
}

type ServerConfig struct {
	Address     string `env:"ADDRESS"`
	FileHandler FileHandlerConfig
}

func (cfg *ServerConfig) FromEnv() *ServerConfig {
	if err := env.Parse(cfg); err != nil {
		log.Fatalln("ServerConfig.FromEnv: can't parse env vars, reason: ", err)
	}
	return cfg
}

func (cfg *ServerConfig) FromFlags() *ServerConfig {
	address := flag.String("a", "127.0.0.1:8080", "server address")
	restore := flag.Bool("r", true, "restore metrics to file")
	storeInterval := flag.Duration("i", 300*time.Second, "store interval")
	storeFile := flag.String("f", "/tmp/devops-metrics-db.json", "json file path to store metrics")
	flag.Parse()
	cfg.Address = *address
	cfg.FileHandler.Restore = *restore
	cfg.FileHandler.Interval = *storeInterval
	cfg.FileHandler.File = *storeFile
	return cfg
}

func NewServerConfig() *ServerConfig {
	cfg := &ServerConfig{
		Address:     "127.0.0.1:8080",
		FileHandler: *newFileStorageConfig(),
	}
	return cfg
}
