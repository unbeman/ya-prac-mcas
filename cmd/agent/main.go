package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/unbeman/ya-prac-mcas/internal/agent"
)

const (
	reportAddr     = "127.0.0.1:8080"
	clientTimeout  = 5 * time.Second
	pollInterval   = 1 * time.Second
	reportInterval = 2 * time.Second
	reportTimeout  = 2 * time.Second
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	defer func() {
		cancel()
		log.Println("Agent cancelled")
	}()
	cm := agent.NewAgentMetrics(reportAddr, clientTimeout, reportTimeout, pollInterval, reportInterval)
	cm.DoWork(ctx)
	<-ctx.Done()
}
