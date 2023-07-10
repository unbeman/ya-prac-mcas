// Package utils describes supportive structs and functions.
package utils

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// TickerPool describes worker pool that do work once in a given time interval.
type TickerPool struct {
	wg   sync.WaitGroup
	done chan struct{}
}

// NewTickerPool creates TickerPool instance.
func NewTickerPool() *TickerPool {
	return &TickerPool{wg: sync.WaitGroup{}, done: make(chan struct{})}
}

// AddTask is method for adding new task to worker group and run it with given time interval.
func (tp *TickerPool) AddTask(ctx context.Context, name string, task func(ctx context.Context), interval time.Duration) {
	tp.wg.Add(1)
	go func() {
		defer tp.wg.Done()
		ticker := time.NewTicker(interval)
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				log.Infof("Task %v stopped", name)
				return
			case <-ticker.C:
				task(ctx)
			}
		}
	}()
}

// Wait wait for workers.
func (tp *TickerPool) Wait() {
	tp.wg.Wait()
}
