package queue

import (
	"context"
	"sync"
)

type Message struct {
	Topic string
	Body  []byte
}
type Publisher interface {
	Publish(context.Context, Message) error
}
type Memory struct {
	mu       sync.Mutex
	Messages []Message
}

func New() *Memory { return &Memory{Messages: []Message{}} }
func (m *Memory) Publish(ctx context.Context, msg Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.Messages = append(m.Messages, msg)
	return nil
}
