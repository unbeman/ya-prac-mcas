package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/agent"
	"github.com/unbeman/ya-prac-mcas/internal/logging"
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

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	go func() {
		exit := make(chan os.Signal, 1)
		signal.Notify(
			exit,
			os.Interrupt,
			syscall.SIGTERM,
			syscall.SIGINT,
			syscall.SIGQUIT,
		)
		<-exit
		cancel()
		log.Println("Agent cancelled")
	}()

	cfg := configs.NewAgentConfig().FromFlags().FromEnv()

	logging.InitLogger(cfg.Logger)

	log.Infof("AGENT CONFIG %+v\n", cfg)

	am, err := agent.NewAgentMetrics(cfg)
	if err != nil {
		log.Error(err)
		return
	}
	am.Run(ctx)
}
