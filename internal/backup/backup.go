package backup

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/parser"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
)

type Backuper interface {
	Restore() error
	Backup() error
	Run()
	Shutdown()
}

type RepositoryBackup struct {
	repo     storage.Repository
	filename string
	interval time.Duration
	closing  chan struct{}
}

func NewRepositoryBackup(cfg *configs.BackupConfig, repo storage.Repository) (*RepositoryBackup, error) {
	rb := &RepositoryBackup{
		filename: cfg.File,
		repo:     repo,
		interval: cfg.Interval,
		closing:  make(chan struct{}),
	}
	if cfg.Restore {
		if err := rb.Restore(); err != nil {
			log.Printf("Can't restore metrics, reason: %v\n", err)
		}
	}
	return rb, nil
}

func (rb *RepositoryBackup) Backup() error {
	log.Debugln("Saving to", rb.filename)
	file, err := os.OpenFile(rb.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return fmt.Errorf("RepositoryBackup.Backup(): can't open file %v - %w", rb.filename, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Errorln(err)
		}
	}()
	writer := json.NewEncoder(file)

	metricsL := rb.repo.GetAll()
	jsonMetricsL := make([]*parser.JSONMetric, 0, len(metricsL))
	for _, metric := range metricsL {
		jsonMetric := parser.MetricToJSON(metric)

		jsonMetricsL = append(jsonMetricsL, jsonMetric)
	}
	log.Debugf("RepositoryBackup.Backup() metrics list for saving %+v\n", jsonMetricsL)

	err = writer.Encode(jsonMetricsL)
	if err != nil {
		return fmt.Errorf("RepositoryBackup.Backup(): %w", err)
	}
	log.Infoln("Metrics saved")
	return nil
}

func (rb *RepositoryBackup) Restore() error {
	file, err := os.OpenFile(rb.filename, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		return fmt.Errorf("RepositoryBackup.Restore(): can't open file %v - %w", rb.filename, err)
	}
	defer file.Close()
	reader := json.NewDecoder(file)
	var jsonMetricsL []*parser.JSONMetric
	err = reader.Decode(&jsonMetricsL)
	if errors.Is(err, io.EOF) {
		log.Infoln("No json metrics to load")
		return nil
	}
	if err != nil { //TODO
		return fmt.Errorf("RepositoryBackup.Restore(): can't decode json %w", err)
	}
	log.Debugf("RepositoryBackup.Restore() metrics %+v\n", jsonMetricsL)
	for _, jsonMetric := range jsonMetricsL {
		params, err := parser.ParseJSON(jsonMetric, parser.PType, parser.PName, parser.PValue)
		if err != nil {
			return fmt.Errorf("RepositoryBackup.Restore(): can't parse json %w", err)
		}
		switch params.Type {
		case metrics.GaugeType:
			rb.repo.SetGauge(params.Name, *params.ValueGauge)
		case metrics.CounterType:
			rb.repo.AddCounter(params.Name, *params.ValueCounter)
		}

	}
	log.Infoln("Metrics loaded")
	return nil
}

func (rb *RepositoryBackup) isTickerEnable() bool {
	if rb.interval != 0*time.Second {
		return true
	}
	return false
}

func (rb *RepositoryBackup) Run() {
	if !rb.isTickerEnable() {
		return
	}
	ticker := time.NewTicker(rb.interval)
	for {
		select {
		case <-rb.closing:
			log.Infoln("Backup ticker stopped")
			return
		case <-ticker.C:
			if err := rb.Backup(); err != nil {
				log.Fatalln(err)
			}
		}
	}
}

func (rb *RepositoryBackup) Shutdown() {
	log.Infoln("Stop backup ticker")
	if rb.isTickerEnable() {
		rb.closing <- struct{}{}
	}
}
