// Package server describes metric server application.
package server

import (
	"context"
	"crypto/rsa"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"sync"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/controller"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
	"github.com/unbeman/ya-prac-mcas/internal/utils"
)

type Server interface {
	GetAddress() string
	Run() error
	Close() error
}

func GetServer(protocol string, addr string, control *controller.Controller, key *rsa.PrivateKey) Server {
	switch protocol {
	case configs.GRPCProtocol:
		return NewGRPCServer(addr, control)
	default:
		return NewHTTPServer(addr, control, key)
	}
}

type application struct {
	repository    storage.Repository
	server        Server
	profileServer *http.Server
}

func GetApplication(cfg configs.ServerConfig) (*application, error) {
	repository, err := storage.GetRepository(cfg.Repository)
	if err != nil {
		return nil, fmt.Errorf("сan't create repository, reason: %w", err)
	}
	privateKey, err := utils.GetPrivateKey(cfg.PrivateCryptoKeyPath)
	if errors.Is(err, utils.ErrNoRSAKey) {
		log.Warning("no private RSA key. Decryption disabled.")
	} else {
		return nil, fmt.Errorf("сan't get private key, reason: %w", err)
	}

	control := controller.NewController(repository, cfg.HashKey)

	server := GetServer(cfg.Protocol, cfg.CollectorAddress, control, privateKey)

	return &application{
		server:        server,
		profileServer: &http.Server{Addr: cfg.ProfileAddress},
		repository:    repository,
	}, nil
}

func (a *application) Start() {
	wg := sync.WaitGroup{}

	// run server
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := a.server.Run()
		log.Info("server closed:", err)
	}()

	// run backup ticker
	if backuper, ok := a.repository.(storage.Backuper); ok {
		wg.Add(1)
		go func() {
			defer wg.Done()
			backuper.Run()
			log.Debugf("Backupper finished")
		}()
	}

	wg.Add(1)
	go func(server *http.Server) {
		defer wg.Done()
		err := server.ListenAndServe()
		log.Infoln("profile server closed:", err)
	}(a.profileServer)

	log.Infoln("Application started, addr:", a.server.GetAddress())

	wg.Wait()

	if backuper, ok := a.repository.(storage.Backuper); ok {
		err := backuper.Backup()
		if err != nil {
			log.Error(err)
		}
	}

	log.Infoln("Application stopped, addr:", a.server.GetAddress())
}

func (a *application) Stop() {
	log.Infoln("Shutting down")
	err := a.server.Close()
	if err != nil {
		log.Error(err)
	}

	err = a.profileServer.Shutdown(context.TODO())
	if err != nil {
		log.Error(err)
	}

	err = a.repository.Shutdown()
	if err != nil {
		log.Error(err)
	}
}
