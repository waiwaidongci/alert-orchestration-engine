package application

import (
	"context"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
	"log/slog"
	"sync"
	"time"
)

type RetryWorker struct {
	log      *slog.Logger
	repo     NotificationRepository
	stop     chan struct{}
	done     chan struct{}
	once     sync.Once
	interval time.Duration
}

func NewRetryWorker(log *slog.Logger, repo NotificationRepository) *RetryWorker {
	return &RetryWorker{log: log, repo: repo, stop: make(chan struct{}), done: make(chan struct{}), interval: 30 * time.Second}
}
func (w *RetryWorker) Start(ctx context.Context) {
	go func() {
		defer close(w.done)
		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-w.stop:
				return
			case <-ticker.C:
				w.tick(ctx)
			}
		}
	}()
}
func (w *RetryWorker) Stop() { w.once.Do(func() { close(w.stop) }); <-w.done }
func (w *RetryWorker) tick(ctx context.Context) {
	items, err := w.repo.List(ctx, "", 100)
	if err != nil {
		w.log.Error("retry list failed", "error", err)
		return
	}
	for _, r := range items {
		if r.Status == notification.Failed && r.CanRetry() {
			w.log.Info("notification eligible for retry", "notification_id", r.ID)
		}
	}
}
