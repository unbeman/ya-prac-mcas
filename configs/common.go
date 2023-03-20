package configs

const (
	ServerAddressDefault = "127.0.0.1:8080"
	KeyDefault           = ""
)

var LogLevelDefault = "info"

type LoggerConfig struct {
	Level string `env:"LOG_LEVEL"`
}

func newLoggerConfig() LoggerConfig {
	return LoggerConfig{Level: LogLevelDefault}
}
