package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"github.com/unbeman/ya-prac-mcas/configs"
	"github.com/unbeman/ya-prac-mcas/internal/metrics"
	"github.com/unbeman/ya-prac-mcas/internal/parser"
	"github.com/unbeman/ya-prac-mcas/internal/storage"
)

type fileHandler struct { // или не хэндлер
	filename string
	interval time.Duration
	repo     storage.Repository
}

func NewFileHandler(cfg configs.FileHandlerConfig, repo storage.Repository) (*fileHandler, error) {
	if len(cfg.File) == 0 {
		return nil, fmt.Errorf("NewFileHandler: no file")
	}
	return &fileHandler{
		filename: cfg.File,
		interval: cfg.Interval,
		repo:     repo,
	}, nil
}

func (fh *fileHandler) Save() error {
	if len(fh.filename) == 0 {
		log.Println("No file for saving")
		return nil
	}
	file, err := os.OpenFile(fh.filename, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return fmt.Errorf("NewFileHandler: can't open file %v - %w", fh.filename, err)
	}
	defer func() {
		if err := file.Close(); err != nil {
			log.Println(err)
		}
	}()
	writer := json.NewEncoder(file)

	metricsL := fh.repo.GetAll()
	jsonMetricsL := make([]*parser.JSONMetric, 0, len(metricsL))
	for _, metric := range metricsL {
		jsonMetric := parser.MetricToJSON(metric)
		jsonMetricsL = append(jsonMetricsL, jsonMetric)
	}
	err = writer.Encode(jsonMetricsL)
	if err != nil {
		return fmt.Errorf("fileHandler.Save(): %w", err)
	}
	log.Println("Metrics saved")
	return nil
}

func (fh *fileHandler) Load() error {
	file, err := os.OpenFile(fh.filename, os.O_RDONLY|os.O_CREATE, 0777)
	if err != nil {
		return fmt.Errorf("NewFileHandler: can't open file %v - %w", fh.filename, err)
	}
	defer file.Close()
	reader := json.NewDecoder(file)
	var jsonMetricsL []*parser.JSONMetric
	err = reader.Decode(&jsonMetricsL)
	if err == io.EOF {
		log.Println("No more json metrics to load")
	}
	if err != nil { //TODO
		return fmt.Errorf("fileHandler.Load(): can't decode json %w", err)
	}
	for _, jsonMetric := range jsonMetricsL {
		params, err := parser.ParseJSON(jsonMetric, parser.PType, parser.PName, parser.PValue)
		if err != nil {
			return fmt.Errorf("fileHandler.Load(): can't parse json %w", err)
		}
		switch params.Type {
		case metrics.GaugeType:
			fh.repo.SetGauge(params.Name, *params.ValueGauge)
		case metrics.CounterType:
			fh.repo.AddCounter(params.Name, *params.ValueCounter)
		}

	}
	log.Println("Metrics loaded")
	return nil
}

func (fh fileHandler) RunSaver(ctx context.Context) {
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
