package configs

import (
	"log"

	"github.com/caarlos0/env/v6"
)

//type ServerConfig interface {
//	Address() string
//}

type ServerConfig struct {
	Address string `env:"ADDRESS"`
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
	}
	return cfg
}
