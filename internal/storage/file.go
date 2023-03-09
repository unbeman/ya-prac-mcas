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
	"github.com/unbeman/ya-prac-mcas/internal/parser"
)

type FileStorage struct { // можт в какой-нибудь Backuper переименовать
	repo     Repository
	filename string
	restore  bool
	interval time.Duration
}

func NewFileStorage(cfg *configs.FileStorageConfig, repo Repository) (*FileStorage, error) {
	if len(cfg.File) == 0 {
		return nil, fmt.Errorf("NewFileStorage: Can't create FileStorage - no file")
	}
	fs := &FileStorage{
		filename: cfg.File,
		interval: cfg.Interval,
		repo:     repo,
	}
	return fs, nil
}

func (fh *FileStorage) Backup() error {
	log.Debugln("Saving to", fh.filename)
	file, err := os.OpenFile(fh.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return fmt.Errorf("NewFileStorage: can't open file %v - %w", fh.filename, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Errorln(err)
		}
	}()
	writer := json.NewEncoder(file)

	metricsL := fh.repo.GetAll()
	jsonMetricsL := make([]*parser.JSONMetric, 0, len(metricsL))
	for _, metric := range metricsL {
		jsonMetric := parser.MetricToJSON(metric)

		jsonMetricsL = append(jsonMetricsL, jsonMetric)
	}
	log.Debugf("FileStorage.Backup() metrics list for saving %+v\n", jsonMetricsL)

	err = writer.Encode(jsonMetricsL)
	if err != nil {
		return fmt.Errorf("FileStorage.Backup(): %w", err)
	}
	log.Infoln("Metrics saved")
	return nil
}

func (fh *FileStorage) Restore() error {
	file, err := os.OpenFile(fh.filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("FileStorage.Restore(): can't open file %v - %w", fh.filename, err)
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
		return fmt.Errorf("FileStorage.Restore(): can't decode json %w", err)
	}
	log.Debugf("FileStorage.Restore() metrics %+v\n", jsonMetricsL)
	for _, jsonMetric := range jsonMetricsL {
		params, err := parser.ParseJSON(jsonMetric, parser.PType, parser.PName, parser.PValue)
		if err != nil {
			return fmt.Errorf("FileStorage.Restore(): can't parse json %w", err)
		}
		switch params.Type {
		case metrics.GaugeType:
			fh.repo.SetGauge(params.Name, *params.ValueGauge)
		case metrics.CounterType:
			fh.repo.AddCounter(params.Name, *params.ValueCounter)
		}

	}
	log.Infoln("Metrics loaded")
	return nil
}

func (fh *FileStorage) RunBackuper(ctx context.Context) {
	if fh.interval == 0*time.Second {
		return
	}
	ticker := time.NewTicker(fh.interval)
	for {
		select {
		case <-ctx.Done():
			log.Println("RunBackuper stopped by context")
			return
		case <-ticker.C:
			if err := fh.Backup(); err != nil {
				log.Fatalln(err)
			}
		}
	}
}
