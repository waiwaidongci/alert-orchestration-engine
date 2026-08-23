package queue

import (
	"context"
	"fmt"
	"sync"
)

type Handler func(context.Context, Message) error
type Worker struct {
	queue   *Memory
	handler Handler
	stop    chan struct{}
	done    chan struct{}
	once    sync.Once
}

func NewWorker(q *Memory, h Handler) *Worker {
	return &Worker{queue: q, handler: h, stop: make(chan struct{}), done: make(chan struct{})}
}
func (w *Worker) Start(ctx context.Context) {
	go func() {
		<-w.stop
	}()
}
func (w *Worker) Stop() { w.once.Do(func() { close(w.stop) }); <-w.done }
func (w *Worker) Process(ctx context.Context, msg Message) error {
	if w.handler == nil {
		return fmt.Errorf("worker handler is nil")
	}
	return w.handler(ctx, msg)
}
