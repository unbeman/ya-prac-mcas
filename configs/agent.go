package configs

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	ServerAddressDefault  = "127.0.0.1:8080"
	PollIntervalDefault   = 2 * time.Second
	ReportIntervalDefault = 10 * time.Second
	ReportTimeoutDefault  = 2 * time.Second
	ClientTimeoutDefault  = 5 * time.Second
)

type AgentOption func(config *AgentConfig)

type HttConnectionConfig struct {
	ClientTimeout time.Duration
}

func newHttConnectionConfig() HttConnectionConfig {
	return HttConnectionConfig{ClientTimeout: ClientTimeoutDefault}
}

type AgentConfig struct {
	Address        string        `env:"ADDRESS"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
	ReportTimeout  time.Duration
	Connection     HttConnectionConfig
	Logger         LoggerConfig
}

func (cfg *AgentConfig) FromEnv() *AgentConfig {
	if err := env.Parse(cfg); err != nil {
		log.Fatalln("AgentConfig.FromEnv: can't parse env vars, reason: ", err)
	}
	return cfg
}

func (cfg *AgentConfig) FromFlags() *AgentConfig {
	address := flag.String("a", ServerAddressDefault, "metrics collection server address")
	pollInterval := flag.Duration("p", PollIntervalDefault, "poll interval")
	reportInterval := flag.Duration("r", ReportIntervalDefault, "report interval")
	logLevel := flag.String("l", LogLevelDefault, "log level, allowed [info, debug]")
	flag.Parse()
	cfg.Address = *address
	cfg.PollInterval = *pollInterval
	cfg.ReportInterval = *reportInterval
	cfg.Logger.Level = *logLevel
	return cfg
}

func NewAgentConfig() *AgentConfig {
	cfg := &AgentConfig{
		Address:        ServerAddressDefault,
		PollInterval:   PollIntervalDefault,
		ReportInterval: ReportIntervalDefault,
		ReportTimeout:  ReportTimeoutDefault,
		Connection:     newHttConnectionConfig(),
		Logger:         newLoggerConfig(),
	}
	return cfg
}
