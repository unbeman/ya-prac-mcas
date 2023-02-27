package configs

import (
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

type HttConnectionConfig struct {
	ClientTimeout time.Duration
}

func newHttConnectionConfig() *HttConnectionConfig {
	return &HttConnectionConfig{ClientTimeout: 5 * time.Second}
}

type AgentConfig struct {
	Address        string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	ReportTimeout  time.Duration
	Connection     *HttConnectionConfig
}

func (cfg *AgentConfig) FromEnv() *AgentConfig {
	if err := env.Parse(cfg); err != nil {
		log.Fatalln("AgentConfig.FromEnv: can't parse env vars, reason: ", err)
	}
	return cfg
}

func NewAgentConfig() *AgentConfig {
	cfg := &AgentConfig{
		Address:        "127.0.0.1:8080",
		PollInterval:   2 * time.Second,
		ReportInterval: 10 * time.Second,
		ReportTimeout:  2 * time.Second,
		Connection:     newHttConnectionConfig(),
	}
	return cfg
}
