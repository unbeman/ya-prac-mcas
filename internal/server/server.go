// Package server describes metric server application.
package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/handlers"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
	"github.com/unbeman/ya-prac-mcas/internal/utils"
)

type serverCollector struct {
	repository    storage.Repository
	httpServer    *http.Server
	profileServer *http.Server
}

func (s *serverCollector) Run() {
	wg := sync.WaitGroup{}

	// run http server
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := s.httpServer.ListenAndServe()
		log.Infoln("collector server stopped ", err)
	}()

	// run backup ticker
	if backuper, ok := s.repository.(storage.Backuper); ok {
		wg.Add(1)
		go func() {
			defer wg.Done()
			backuper.Run()
			log.Debugf("Backupper Run is end")
		}()
	}

	wg.Add(1)
	go func(server *http.Server) {
		defer wg.Done()
		err := server.ListenAndServe()
		log.Infoln("profile server stopped:", err)
	}(s.profileServer)

	log.Infoln("Application started, addr:", s.httpServer.Addr)

	wg.Wait()

	if backuper, ok := s.repository.(storage.Backuper); ok {
		err := backuper.Backup()
		if err != nil {
			log.Error(err)
		}
	}

	log.Infoln("Application stopped, addr:", s.httpServer.Addr)
}

func (s *serverCollector) Shutdown() {
	log.Infoln("Shutting down")
	err := s.httpServer.Shutdown(context.TODO())
	if err != nil {
		log.Errorln(err)
	}

	err = s.profileServer.Shutdown(context.TODO())
	if err != nil {
		log.Errorln(err)
	}

	err = s.repository.Shutdown()
	if err != nil {
		log.Errorln(err)
	}
}

func NewServerCollector(cfg configs.ServerConfig) (*serverCollector, error) {
	repository, err := storage.GetRepository(cfg.Repository)
	if err != nil {
		return nil, fmt.Errorf("сan't create repository, reason: %w", err)
	}
	privateKey, err := utils.GetPrivateKey(cfg.PrivateCryptoKeyPath)
	if errors.Is(err, utils.NoRSAKeyErr) {
		log.Warning("no private RSA key. Decryption disabled.")
	}
	if err != nil {
		return nil, fmt.Errorf("сan't get private key, reason: %w", err)
	}

	handler := handlers.NewCollectorHandler(repository, cfg.HashKey, privateKey)

	return &serverCollector{
		httpServer:    &http.Server{Addr: cfg.CollectorAddress, Handler: handler},
		profileServer: &http.Server{Addr: cfg.ProfileAddress},
		repository:    repository,
	}, nil
}
