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

type Backuper interface {
	Restore() error
	Backup() error
	Run()
	Shutdown() error
}

type BackupRepository struct {
	Repository
	filename string
	interval time.Duration
	closing  chan struct{}
}

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

func (br *BackupRepository) Backup() error {
	log.Debugln("Saving to", br.filename)
	file, err := os.OpenFile(br.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return fmt.Errorf("BackupRepository.Backup(): can't open file %v - %w", br.filename, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Errorln(err)
		}
	}()
	writer := json.NewEncoder(file)

	metricsL, err := br.GetAll(context.TODO())
	if err != nil {
		return fmt.Errorf("BackupRepository.Backup(): %v", err)
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
	log.Infoln("Metrics saved")
	return nil
}

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
		log.Infoln("No json metrics to load")
		return nil
	}
	if err != nil { //TODO
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
			return fmt.Errorf("BackupRepository.Restore(): %v", err)
		}

	}
	log.Infoln("Metrics loaded")
	return nil
}

func (br *BackupRepository) isTickerEnable() bool {
	return br.interval != 0*time.Second
}

func (br *BackupRepository) Run() {
	if !br.isTickerEnable() {
		return
	}
	ticker := time.NewTicker(br.interval)
	for {
		select {
		case <-br.closing:
			ticker.Stop()
			log.Infoln("Backup ticker stopped")
			return
		case <-ticker.C:
			if err := br.Backup(); err != nil {
				log.Error(err)
			}
		}
	}
}

func (br *BackupRepository) Shutdown() error {
	log.Infoln("Stop backup ticker")
	if br.isTickerEnable() {
		br.closing <- struct{}{}
	}
	if err := br.Backup(); err != nil {
		log.Error(err)
	}
	return nil
}
