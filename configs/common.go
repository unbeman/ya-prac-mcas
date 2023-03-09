package configs

var LogLevelDefault = "info"

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL"`
}

func newLoggerConfig() LoggerConfig {
	return LoggerConfig{Level: LogLevelDefault}
}
