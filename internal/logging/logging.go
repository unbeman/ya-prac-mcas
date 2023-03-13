package logging

import (
	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/configs"
)

const (
	LogDebug = "debug"
	LogInfo  = "info"
)

func InitLogger(cfg configs.LoggerConfig) {
	switch cfg.Level {
	case LogInfo:
		log.SetLevel(log.InfoLevel)
	case LogDebug:
		log.SetLevel(log.DebugLevel)
	}
}
