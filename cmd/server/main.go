package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-mcas/internal/logging"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/server"
)

func initContext() (context.Context, context.CancelFunc) {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	return ctx, cancel
}

func initConfig() configs.ServerConfig {
	return *configs.NewServerConfig(configs.FromFlags(), configs.FromEnv())
}

// TODO: wrap to init server
func main() { //TODO: more logs, pass context to Repository methods and handlers
	ctx, cancel := initContext()
	defer func() {
		cancel()
		log.Println("Server cancelled")
	}()

	cfg := initConfig()

	logging.InitLogger(cfg.Logger)

	log.Debugf("SERVER CONFIG %+v\n", cfg)

	collectorServer := server.NewServerCollector(cfg)
	collectorServer.Run(ctx)
	<-ctx.Done()
	collectorServer.Shutdown()
}
