package utils

import (
	"context"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

type TickerPool struct {
	wg   sync.WaitGroup
	done chan struct{}
}

func NewTickerPool() *TickerPool {
	return &TickerPool{wg: sync.WaitGroup{}, done: make(chan struct{})}
}

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

func (tp *TickerPool) Wait() {
	tp.wg.Wait()
}
