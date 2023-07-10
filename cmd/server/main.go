package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/internal/logging"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/server"
)

// Build info variables
//
// Example:
//
// go run -ldflags "-X 'main.buildVersion=v1.0.0' -X 'main.buildDate=$(date +'%Y/%m/%d %H:%M:%S')' -X 'main.buildCommit=$(git rev-parse HEAD)'" cmd/agent/main.go
var (
	buildVersion = "N/A"
	buildDate    = "N/A"
	buildCommit  = "N/A"
)

func main() {
	fmt.Printf("Build version: %v\n", buildVersion)
	fmt.Printf("Build date: %v\n", buildDate)
	fmt.Printf("Build commit: %v\n", buildCommit)

	cfg := *configs.NewServerConfig(configs.FromFlags(), configs.FromEnv())

	logging.InitLogger(cfg.Logger)

	log.Debugf("SERVER CONFIG %+v\n", cfg)

	collectorServer, err := server.NewServerCollector(cfg)
	if err != nil {
		log.Error(err)
		return
	}

	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(
			exit,
			os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGQUIT,
		)

		sig := <-exit
		log.Infof("Got signal '%v'", sig)

		collectorServer.Shutdown()
	}()

	collectorServer.Run()
}
