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

// TODO: wrap to init agent
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	defer func() {
		cancel()
		log.Println("Agent cancelled")
	}()

	cfg := configs.NewAgentConfig().FromFlags().FromEnv()

	logging.InitLogger(cfg.Logger)

	log.Infof("AGENT CONFIG %+v\n", cfg)

	cm := agent.NewAgentMetrics(cfg)
	cm.DoWork(ctx)
	<-ctx.Done()
}
