package server

import (
	"context"
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/handlers"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
)

type serverCollector struct {
	address     string
	handler     http.Handler
	repository  storage.Repository
	fileStorage *storage.FileStorage
	restore     bool
}

func (s serverCollector) Run(ctx context.Context) {
	if s.fileStorage != nil {
		if s.restore {
			if err := s.fileStorage.Restore(); err != nil {
				log.Println("Can't restore metrics, reason: %v\n", err)
			}
		}

		go s.fileStorage.RunBackuper(ctx)
	}

	go func(ctx context.Context) {
		log.Fatalln(http.ListenAndServe(s.address, s.handler))
	}(ctx)
	log.Infoln("Server started, addr:", s.address)
}

func (s serverCollector) Shutdown() {
	if s.fileStorage == nil {
		return
	}
	err := s.fileStorage.Backup()
	if err != nil {
		log.Fatalln("Can't backup metrics, reason:", err)
	}
}

func NewServerCollector(cfg configs.ServerConfig) *serverCollector {
	repository := storage.GetRepository()
	fileStorage, err := storage.NewFileStorage(cfg.FileStorage, repository)
	if err != nil {
		log.Warnf("FileStorage disabled, %v", err)
	}
	return &serverCollector{
		address:     cfg.Address,
		handler:     handlers.NewCollectorHandler(repository),
		repository:  repository,
		restore:     cfg.Restore,
		fileStorage: fileStorage,
	}
}
