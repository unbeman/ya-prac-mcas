package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/agent"
)

// TODO: wrap to init agent
func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	defer func() {
		cancel()
		log.Println("Agent cancelled")
	}()

	cfg := configs.NewAgentConfig().FromEnv()

	client := http.Client{Timeout: cfg.Connection.ClientTimeout}

	cm := agent.NewAgentMetrics(cfg, &client)
	cm.DoWork(ctx)
	<-ctx.Done()
}
