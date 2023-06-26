package configs

import (
	"flag"
	"log"
	"time"

	"github.com/caarlos0/env/v6"
)

const (
	PollIntervalDefault    = 2 * time.Second
	ReportIntervalDefault  = 10 * time.Second
	ReportTimeoutDefault   = 2 * time.Second
	ClientTimeoutDefault   = 5 * time.Second
	RateTokensCountDefault = 100
)

type AgentOption func(config *AgentConfig)

type HttConnectionConfig struct {
	Address         string `env:"ADDRESS"`
	RateTokensCount int
	ClientTimeout   time.Duration
	ReportTimeout   time.Duration
}

func newHttConnectionConfig() HttConnectionConfig {
	return HttConnectionConfig{
		Address:         ServerAddressDefault,
		ClientTimeout:   ClientTimeoutDefault,
		ReportTimeout:   ReportTimeoutDefault,
		RateTokensCount: RateTokensCountDefault,
	}
}

type AgentConfig struct {
	Key            string        `env:"KEY"`
	PollInterval   time.Duration `env:"POLL_INTERVAL"`
	ReportInterval time.Duration `env:"REPORT_INTERVAL"`
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
	rateTokensCount := flag.Int("l", RateTokensCountDefault, "limit request count in one second")
	key := flag.String("k", KeyDefault, "key for calculating the metric hash")
	pollInterval := flag.Duration("p", PollIntervalDefault, "poll interval")
	reportInterval := flag.Duration("r", ReportIntervalDefault, "report interval")
	logLevel := flag.String("e", LogLevelDefault, "log level, allowed [info, debug]")
	flag.Parse()
	cfg.Key = *key
	cfg.PollInterval = *pollInterval
	cfg.ReportInterval = *reportInterval
	cfg.Connection.Address = *address
	cfg.Connection.RateTokensCount = *rateTokensCount
	cfg.Logger.Level = *logLevel
	return cfg
}

func NewAgentConfig() *AgentConfig {
	cfg := &AgentConfig{
		Key:            KeyDefault,
		PollInterval:   PollIntervalDefault,
		ReportInterval: ReportIntervalDefault,
		Connection:     newHttConnectionConfig(),
		Logger:         newLoggerConfig(),
	}
	return cfg
}
