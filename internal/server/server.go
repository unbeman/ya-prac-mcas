package server

import (
	"context"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/handlers"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
)

type serverCollector struct {
	repository storage.Repository
	httpServer http.Server
}

func (s *serverCollector) Run() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	// run http server
	go func() {
		defer wg.Done()
		err := s.httpServer.ListenAndServe()
		log.Infoln("serverCollector.Run()", err)
	}()

	log.Infoln("Server collector started, addr:", s.httpServer.Addr)

	wg.Wait()

	log.Infoln("Server collector stopped, addr:", s.httpServer.Addr)
}

func (s *serverCollector) Shutdown() {

	log.Infoln("Shutting down")
	err := s.httpServer.Shutdown(context.TODO())
	if err != nil {
		log.Errorln(err)
	}
	//не успевает
	err = s.repository.Shutdown()
	if err != nil {
		log.Errorln(err)
	}
}

func NewServerCollector(cfg configs.ServerConfig) *serverCollector {
	repository, err := storage.GetRepository(cfg.Repository)
	if err != nil {
		log.Fatalln("Can't create repository, reason:", err)
	}

	handler := handlers.NewCollectorHandler(repository, cfg.Key)
	return &serverCollector{
		httpServer: http.Server{Addr: cfg.Address, Handler: handler},
		repository: repository,
	}
}
