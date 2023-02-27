package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/handlers"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
)

// TODO: wrap to init server
func main() { //TODO: more logs, add context to Repository and handlers
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT, os.Interrupt)
	defer func() {
		cancel()
		log.Println("Server cancelled")
	}()
	cfg := configs.NewServerConfig().FromFlags().FromEnv()
	ramRepo := storage.NewRAMRepository()

	fileHandler, err := handlers.NewFileHandler(cfg.FileHandler, ramRepo)
	if err != nil {
		log.Println("Can't init file storage, skipped, reason:", err)
	}
	if fileHandler != nil {
		defer func() {
			if err := fileHandler.Save(); err != nil {
				log.Fatalln(err)
			}
		}()

		if cfg.FileHandler.Restore {
			err := fileHandler.Load()
			if err != nil {
				log.Println("Can't restore RAMRepository, skipped, reason:", err)
			}
		}

		go fileHandler.RunSaver(ctx)
	}

	collectorHandler := handlers.NewCollectorHandler(ramRepo)

	go func(ctx context.Context) {
		log.Fatal(http.ListenAndServe(cfg.Address, collectorHandler))
	}(ctx)
	log.Println("Server started, addr:", cfg.Address)
	<-ctx.Done()
}
