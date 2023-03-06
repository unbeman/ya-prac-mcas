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

type fileRepository struct {
	ramRepository
	filename string
	restore  bool
	interval time.Duration
}

func NewFileRepository(cfg configs.FileStorageConfig) *fileRepository {
	fs := &fileRepository{
		filename:      cfg.File,
		restore:       cfg.Restore,
		interval:      cfg.Interval,
		ramRepository: *NewRAMRepository(),
	}
	return fs
}

func (fh *fileRepository) Save() error {
	log.Debugln("Saving to", fh.filename)
	file, err := os.OpenFile(fh.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
	if err != nil {
		return fmt.Errorf("NewFileRepository: can't open file %v - %w", fh.filename, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Errorln(err)
		}
	}()
	writer := json.NewEncoder(file)

	metricsL := fh.GetAll()
	jsonMetricsL := make([]*parser.JSONMetric, 0, len(metricsL))
	for _, metric := range metricsL {
		jsonMetric := parser.MetricToJSON(metric)

		jsonMetricsL = append(jsonMetricsL, jsonMetric)
	}
	log.Debugf("fileRepository.Save() metrics list for saving %+v\n", jsonMetricsL)

	err = writer.Encode(jsonMetricsL)
	if err != nil {
		return fmt.Errorf("fileRepository.Save(): %w", err)
	}
	log.Infoln("Metrics saved")
	return nil
}

func (fh *fileRepository) Load() error {
	if !fh.restore {
		return nil
	}
	file, err := os.OpenFile(fh.filename, os.O_RDONLY|os.O_CREATE, 0664)
	if err != nil {
		return fmt.Errorf("NewFileRepository: can't open file %v - %w", fh.filename, err)
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
		return fmt.Errorf("fileRepository.Load(): can't decode json %w", err)
	}
	log.Debugf("fileRepository.Load() metrics %+v\n", jsonMetricsL)
	for _, jsonMetric := range jsonMetricsL {
		params, err := parser.ParseJSON(jsonMetric, parser.PType, parser.PName, parser.PValue)
		if err != nil {
			return fmt.Errorf("fileRepository.Load(): can't parse json %w", err)
		}
		switch params.Type {
		case metrics.GaugeType:
			fh.SetGauge(params.Name, *params.ValueGauge)
		case metrics.CounterType:
			fh.AddCounter(params.Name, *params.ValueCounter)
		}

	}
	log.Infoln("Metrics loaded")
	return nil
}

func (fh *fileRepository) RunSaver(ctx context.Context) {
	if fh.interval == 0*time.Second {
		return
	}
	ticker := time.NewTicker(fh.interval)
	for {
		select {
		case <-ctx.Done():
			log.Println("RunSaver stopped by context")
			return
		case <-ticker.C:
			if err := fh.Save(); err != nil {
				log.Fatalln(err)
			}
		}
	}
}
