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
	address    string
	handler    http.Handler
	repository storage.Repository
}

func (s serverCollector) Run(ctx context.Context) {
	err := s.repository.Load()
	if err != nil {
		log.Fatalln("Can't load metrics")
	}
	go s.repository.RunSaver(ctx)

	defer func() {
		err := s.repository.Save()
		if err != nil {
			log.Fatalln("Can't save metrics, reason:", err)
		}
	}()
	go func(ctx context.Context) {
		log.Fatalln(http.ListenAndServe(s.address, s.handler))
	}(ctx)
	log.Infoln("Server started, addr:", s.address)
	<-ctx.Done()
}
func NewServerCollector(cfg configs.ServerConfig) *serverCollector {
	repository := storage.GetRepository(cfg.Repository)
	collectorHandler := handlers.NewCollectorHandler(repository)
	return &serverCollector{
		address:    cfg.Address,
		handler:    collectorHandler,
		repository: repository,
	}
}
