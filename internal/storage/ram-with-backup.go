package storage

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
)

// Backuper describes interface for saving metrics to file and loading ones to memory.
type Backuper interface {
	Restore() error
	Backup() error
	Run()
	Shutdown() error
}

// BackupRepository is implementation of Backuper and Repository.
type BackupRepository struct {
	Repository
	filename string
	interval time.Duration
	closing  chan struct{}
}

// NewRAMBackupRepository initialize new BackupRepository with config.
func NewRAMBackupRepository(cfg *configs.BackupConfig) (*BackupRepository, error) {
	if len(cfg.File) == 0 {
		return nil, errors.New("no filename")
	}
	rb := &BackupRepository{
		filename:   cfg.File,
		Repository: NewRAMRepository(),
		interval:   cfg.Interval,
		closing:    make(chan struct{}),
	}
	if cfg.Restore {
		if err := rb.Restore(); err != nil {
			log.Printf("Can't restore metrics, reason: %v\n", err)
		}
	}

	return rb, nil
}

// Backup saves metrics from memory storage to file.
func (br *BackupRepository) Backup() error {
	log.Debug("Saving to", br.filename)
	file, err := os.OpenFile(br.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return fmt.Errorf("BackupRepository.Backup(): can't open file %v - %w", br.filename, err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Errorln(err)
		}
	}()
	writer := json.NewEncoder(file)

	metricsL, err := br.GetAll(context.TODO())
	if err != nil {
		return fmt.Errorf("BackupRepository.Backup(): %w", err)
	}
	jsonMetricsL := make([]*metrics.Params, 0, len(metricsL))
	for _, metric := range metricsL {
		jsonMetric := metric.ToParams()
		jsonMetricsL = append(jsonMetricsL, &jsonMetric)
	}
	log.Debugf("BackupRepository.Backup() metrics list for saving %+v\n", jsonMetricsL)

	err = writer.Encode(jsonMetricsL)
	if err != nil {
		return fmt.Errorf("BackupRepository.Backup(): %w", err)
	}
	log.Info("Metrics saved")
	return nil
}

// Restore loads metrics from file to memory storage.
func (br *BackupRepository) Restore() error {
	file, err := os.OpenFile(br.filename, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		return fmt.Errorf("BackupRepository.Restore(): can't open file %v - %w", br.filename, err)
	}
	defer file.Close()
	reader := json.NewDecoder(file)
	var jsonMetricsL []*metrics.Params
	err = reader.Decode(&jsonMetricsL)
	if errors.Is(err, io.EOF) {
		log.Info("No json metrics to load")
		return nil
	}
	if err != nil {
		return fmt.Errorf("BackupRepository.Restore(): can't decode json %w", err)
	}
	log.Debugf("BackupRepository.Restore() metrics %+v\n", jsonMetricsL)
	for _, params := range jsonMetricsL {
		switch params.Type {
		case metrics.GaugeType:
			_, err = br.SetGauge(context.TODO(), params.Name, *params.ValueGauge)
		case metrics.CounterType:
			_, err = br.AddCounter(context.TODO(), params.Name, *params.ValueCounter)
		}
		if err != nil {
			return fmt.Errorf("BackupRepository.Restore(): %w", err)
		}

	}
	log.Info("Metrics loaded")
	return nil
}

// isTickerEnable defines need to turn on ticker.
func (br *BackupRepository) isTickerEnable() bool {
	return br.interval != 0*time.Second
}

// Run makes backup every interval, if interval more than 0 seconds.
func (br *BackupRepository) Run() {
	if !br.isTickerEnable() {
		log.Info("Backupper not started, no interval provided")
		return
	}
	log.Info("Backupper started")
	ticker := time.NewTicker(br.interval)
	for {
		select {
		case <-br.closing:
			ticker.Stop()
			log.Info("Backup ticker stopped")
			return
		case <-ticker.C:
			if err := br.Backup(); err != nil {
				log.Error(err)
			}
		}
	}
}

// Shutdown sends signal for stopping interval backup.
func (br *BackupRepository) Shutdown() error {
	log.Info("Stop backup ticker")
	if br.isTickerEnable() {
		br.closing <- struct{}{}
	}
	return nil
}
