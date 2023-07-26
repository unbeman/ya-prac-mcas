// Package configs describes applications settings.
package configs

// Default config settings
const (
	ServerAddressDefault = "127.0.0.1:8080"
	KeyDefault           = ""
)

var LogLevelDefault = "info"

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL" json:"log_level"`
}

func newLoggerConfig() LoggerConfig {
	return LoggerConfig{Level: LogLevelDefault}
}
