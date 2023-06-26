package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/agent"
	"github.com/unbeman/ya-prac-mcas/internal/logging"
)

func main() {
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

	am := agent.NewAgentMetrics(cfg)
	am.Run(ctx)
}
