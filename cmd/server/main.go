package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/internal/logging"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/server"
)

func main() {
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
