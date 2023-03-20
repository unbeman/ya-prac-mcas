package server

import (
	"context"
	"net/http"
	"sync"

	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-mcas/internal/backup"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/handlers"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
)

type serverCollector struct {
	repository storage.Repository
	httpServer http.Server
	backuper   backup.Backuper
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
	if s.isBackuperEnabled() {
		wg.Add(1)
		go func() {
			defer wg.Done()
			s.backuper.Run()
		}()
	}

	log.Infoln("Server collector started, addr:", s.httpServer.Addr)

	wg.Wait()

	// shutdown
	if s.isBackuperEnabled() {
		if err := s.backuper.Backup(); err != nil {
			log.Errorf("Failed to backup repository: %v", err)
		}
	}

	log.Infoln("Server collector stopped, addr:", s.httpServer.Addr)
}

func (s *serverCollector) Shutdown() {

	log.Infoln("Shutting down")
	err := s.httpServer.Shutdown(context.TODO())
	if err != nil {
		log.Errorln(err)
	}
	if s.isBackuperEnabled() {
		s.backuper.Shutdown()
	}
	err = s.repository.Shutdown()
	if err != nil {
		log.Errorln(err)
	}
}

func NewServerCollector(cfg configs.ServerConfig) *serverCollector {
	var (
		backuper backup.Backuper
		err      error
	)

	repository, err := storage.GetRepository(cfg.Repository)
	if err != nil {
		log.Fatalln("Can't create repository, reason:", err)
	}
	if cfg.Repository.PG == nil { //TODO: сделать адекватно
		backuper, err = backup.NewRepositoryBackup(cfg.Backup, repository) //TODO: вернуть FileRepository
		if err != nil {
			log.Fatalln("NewServerCollector:", err)
		}
	}

	handler := handlers.NewCollectorHandler(repository, cfg.Key)
	return &serverCollector{
		httpServer: http.Server{Addr: cfg.Address, Handler: handler},
		repository: repository,
		backuper:   backuper,
	}
}

func (s *serverCollector) isBackuperEnabled() bool {
	return s.backuper != nil
}
