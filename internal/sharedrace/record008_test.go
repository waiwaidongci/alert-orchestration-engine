package r008

import (
	"context"
	httpadapter "github.com/example/alert-orchestration-engine/internal/adapter/http"
	"github.com/example/alert-orchestration-engine/internal/domain/alert"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
	"github.com/example/alert-orchestration-engine/internal/infrastructure/memory"
	"github.com/example/alert-orchestration-engine/internal/infrastructure/notifier"
	"sync"
	"testing"
	"time"
)

func TestR008Metrics(t *testing.T) {
	m := &httpadapter.Metrics{}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; m.Request(); m.Event() }()
	}
	close(start)
	wg.Wait()
}
func TestR008Dispatcher(t *testing.T) {
	d := notifier.New()
	r := notification.New("n", "a", "r", "webhook", "x", time.Now())
	a := &alert.Alert{ID: "a"}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; _ = d.Send(context.Background(), r, a) }()
	}
	close(start)
	wg.Wait()
}
func TestR008IDs(t *testing.T) {
	ids := &memory.IDs{}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() { defer wg.Done(); <-start; _ = ids.NewID("x") }()
	}
	close(start)
	wg.Wait()
}
func TestR008Health(t *testing.T) {
	h := &httpadapter.HealthState{}
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(i int) { defer wg.Done(); <-start; h.SetReady(i%2 == 0); _ = h.Ready() }(i)
	}
	close(start)
	wg.Wait()
}
