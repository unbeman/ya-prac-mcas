package server

import (
	"context"
	"fmt"
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

	// run backup ticker
	if backuper, ok := s.repository.(storage.Backuper); ok {
		wg.Add(1)
		go func() {
			defer wg.Done()
			backuper.Run()
		}()
	}

	log.Infoln("Server collector started, addr:", s.httpServer.Addr)

	wg.Wait()

	log.Infoln("Server collector stopped, addr:", s.httpServer.Addr)

	if backuper, ok := s.repository.(storage.Backuper); ok {
		err := backuper.Backup()
		if err != nil {
			log.Error(err)
		}
	}
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

func NewServerCollector(cfg configs.ServerConfig) (*serverCollector, error) {
	repository, err := storage.GetRepository(cfg.Repository)
	if err != nil {
		return nil, fmt.Errorf("сan't create repository, reason: %v", err)
	}

	handler := handlers.NewCollectorHandler(repository, cfg.Key)
	return &serverCollector{
		httpServer: http.Server{Addr: cfg.Address, Handler: handler},
		repository: repository,
	}, nil
}
