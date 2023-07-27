// Package logging describes logger and it's initialize.
package logging

import (
	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/configs"
)

// Default log settings.
const (
	LogDebug = "debug"
	LogInfo  = "info"
)

// InitLogger set log parameters depending on the input config.
func InitLogger(cfg configs.LoggerConfig) {
	switch cfg.Level {
	case LogInfo:
		log.SetLevel(log.InfoLevel)
	case LogDebug:
		log.SetLevel(log.DebugLevel)
	}
}
