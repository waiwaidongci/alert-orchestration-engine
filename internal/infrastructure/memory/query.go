package memory

import (
	"context"
	"github.com/example/alert-orchestration-engine/internal/domain/alert"
	"github.com/example/alert-orchestration-engine/internal/domain/event"
	"time"
)

func (s *Store) FindAlerts(ctx context.Context, f alert.Filter, limit int) ([]*alert.Alert, error) {
	items, err := s.List(ctx, f.Status, limit*2)
	if err != nil {
		return nil, err
	}
	out := []*alert.Alert{}
	for _, a := range items {
		if a.Matches(f) {
			out = append(out, a)
		}
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}
func (s *Store) LatestEvent(ctx context.Context) (event.Event, error) {
	items, err := s.ListEvents(ctx, 1)
	if err != nil || len(items) == 0 {
		return event.Event{}, err
	}
	return items[0], nil
}
func (s *Store) EventsSince(ctx context.Context, since time.Time) ([]event.Event, error) {
	items, err := s.ListEvents(ctx, 0)
	if err != nil {
		return nil, err
	}
	out := []event.Event{}
	for _, e := range items {
		if !e.ReceivedAt.Before(since) {
			out = append(out, e)
		}
	}
	return out, nil
}
