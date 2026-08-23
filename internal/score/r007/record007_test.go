package r007

import (
	"context"
	"github.com/example/alert-orchestration-engine/internal/application"
	"github.com/example/alert-orchestration-engine/internal/infrastructure/queue"
	"log/slog"
	"testing"
	"time"
)

func waitStops(t *testing.T, stop func(), label string) {
	t.Helper()
	start := make(chan struct{})
	done := make(chan struct{}, 2)
	for i := 0; i < 2; i++ {
		go func() { <-start; stop(); done <- struct{}{} }()
	}
	close(start)
	for i := 0; i < 2; i++ {
		select {
		case <-done:
		case <-time.After(300 * time.Millisecond):
			t.Fatal(label)
		}
	}
}

func TestR007RetryCancel(t *testing.T) {
	w := application.NewRetryWorker(slog.Default(), nil)
	ctx, c := context.WithCancel(context.Background())
	c()
	w.Start(ctx)
	waitStops(t, w.Stop, "retry stop blocked")
}
func TestR007RetryStop(t *testing.T) {
	w := application.NewRetryWorker(slog.Default(), nil)
	w.Start(context.Background())
	waitStops(t, w.Stop, "retry stop blocked")
}
func TestR007QueueCancel(t *testing.T) {
	w := queue.NewWorker(queue.New(), nil)
	ctx, c := context.WithCancel(context.Background())
	c()
	w.Start(ctx)
	waitStops(t, w.Stop, "queue stop blocked")
}
func TestR007QueueStop(t *testing.T) {
	w := queue.NewWorker(queue.New(), nil)
	w.Start(context.Background())
	waitStops(t, w.Stop, "queue stop blocked")
}
