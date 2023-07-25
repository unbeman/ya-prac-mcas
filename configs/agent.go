// Package configs describes applications settings.
package configs

import (
	"encoding/json"
	"flag"
	"log"
	"os"
	"time"

	"github.com/caarlos0/env/v6"
)

// Default config settings
const (
	PollIntervalDefault        = 2 * time.Second
	ReportIntervalDefault      = 10 * time.Second
	ReportTimeoutDefault       = 2 * time.Second
	ClientTimeoutDefault       = 5 * time.Second
	RateTokensCountDefault     = 100
	PublicCryptoKeyPathDefault = "public.pem"
	JSONAgentConfigPathDefault = ""
)

type AgentOption func(config *AgentConfig)

type HttConnectionConfig struct {
	Address         string `env:"ADDRESS" json:"address,omitempty"`
	RateTokensCount int    `json:"rate_tokens_count,omitempty"`
	ClientTimeout   time.Duration
	ReportTimeout   time.Duration
}

func (cfg *HttConnectionConfig) UnmarshalJSON(data []byte) error {
	type RealCfg HttConnectionConfig
	jCfg := struct {
		ClientTimeout string `json:"client_timeout,omitempty"`
		ReportTimeout string `json:"report_timeout,omitempty"`
		*RealCfg
	}{
		RealCfg: (*RealCfg)(cfg),
	}

	err := json.Unmarshal(data, &jCfg)
	if err != nil {
		return err
	}
	if jCfg.ClientTimeout != "" {
		cfg.ClientTimeout, err = time.ParseDuration(jCfg.ClientTimeout)
		if err != nil {
			return err
		}
	}

	if jCfg.ClientTimeout != "" {
		cfg.ReportTimeout, err = time.ParseDuration(jCfg.ReportTimeout)
		if err != nil {
			return err
		}
	}

	return nil
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
	HashKey             string        `env:"KEY" json:"key,omitempty"`
	PublicCryptoKeyPath string        `env:"CRYPTO_KEY" json:"crypto_key,omitempty"`
	PollInterval        time.Duration `env:"POLL_INTERVAL"`
	ReportInterval      time.Duration `env:"REPORT_INTERVAL"`
	Connection          HttConnectionConfig
	Logger              LoggerConfig
	jsonConfigPath      string
}

func (cfg *AgentConfig) UnmarshalJSON(data []byte) error {
	type RealCfg AgentConfig
	jCfg := struct {
		PollInterval   string `json:"poll_interval,omitempty"`
		ReportInterval string `json:"report_interval,omitempty"`
		*RealCfg
	}{
		RealCfg: (*RealCfg)(cfg),
	}

	err := json.Unmarshal(data, &jCfg)
	if err != nil {
		return err
	}
	if jCfg.PollInterval != "" {
		cfg.PollInterval, err = time.ParseDuration(jCfg.PollInterval)
		if err != nil {
			return err
		}
	}

	if jCfg.ReportInterval != "" {
		cfg.ReportInterval, err = time.ParseDuration(jCfg.ReportInterval)
		if err != nil {
			return err
		}
	}

	return nil
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
	publicCryptoKeyPath := flag.String("crypto-key", PublicCryptoKeyPathDefault, "path to public crypto key file")
	pollInterval := flag.Duration("p", PollIntervalDefault, "poll interval")
	reportInterval := flag.Duration("r", ReportIntervalDefault, "report interval")
	logLevel := flag.String("e", LogLevelDefault, "log level, allowed [info, debug]")
	flag.Parse()
	cfg.HashKey = *key

	cfg.PollInterval = *pollInterval
	cfg.ReportInterval = *reportInterval
	cfg.Connection.Address = *address
	cfg.Connection.RateTokensCount = *rateTokensCount
	cfg.PublicCryptoKeyPath = *publicCryptoKeyPath
	cfg.Logger.Level = *logLevel
	return cfg
}

func (cfg *AgentConfig) FromJSON() *AgentConfig {
	jsonConfigPath := flag.String("c", cfg.jsonConfigPath, "path to json config")
	cfg.jsonConfigPath = *jsonConfigPath
	flag.Parse()

	envPath, isSet := os.LookupEnv("CONFIG")
	if isSet {
		cfg.jsonConfigPath = envPath
	}

	if cfg.jsonConfigPath == "" {
		return cfg
	}

	data, err := os.ReadFile(cfg.jsonConfigPath)
	if err != nil {
		log.Fatalf("can't read %v, reason: %v", cfg.jsonConfigPath, err)
	}

	err = json.Unmarshal(data, &cfg)
	if err != nil {
		log.Fatalf("can't unmarshal json config, reason: %v", err)
	}

	err = json.Unmarshal(data, &cfg.Logger)
	if err != nil {
		log.Fatalf("can't unmarshal json config logger, reason: %v", err)
	}

	err = json.Unmarshal(data, &cfg.Connection)
	if err != nil {
		log.Fatalf("can't unmarshal json config connection, reason: %v", err)
	}
	return cfg
}

func NewAgentConfig() *AgentConfig {
	cfg := &AgentConfig{
		HashKey:             KeyDefault,
		PublicCryptoKeyPath: PublicCryptoKeyPathDefault,
		PollInterval:        PollIntervalDefault,
		ReportInterval:      ReportIntervalDefault,
		Connection:          newHttConnectionConfig(),
		Logger:              newLoggerConfig(),
		jsonConfigPath:      JSONAgentConfigPathDefault,
	}
	return cfg
}
