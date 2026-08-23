package memory

import (
	"context"
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/domain/alert"
	"github.com/example/alert-orchestration-engine/internal/domain/event"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
	"github.com/example/alert-orchestration-engine/internal/domain/rule"
	"github.com/example/alert-orchestration-engine/internal/domain/silence"
	"sort"
	"sync"
	"time"
)

type Store struct {
	mu            sync.RWMutex
	events        []event.Event
	alerts        map[string]*alert.Alert
	byFingerprint map[string]string
	rules         map[string]rule.Rule
	silences      map[string]silence.Silence
	notifications map[string]notification.Record
}

func NewStore() *Store {
	return &Store{alerts: map[string]*alert.Alert{}, byFingerprint: map[string]string{}, rules: map[string]rule.Rule{}, silences: map[string]silence.Silence{}, notifications: map[string]notification.Record{}}
}
func (s *Store) SaveEvent(ctx context.Context, e event.Event) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, e)
	return nil
}
func (s *Store) ListEvents(ctx context.Context, limit int) ([]event.Event, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if limit <= 0 || limit > len(s.events) {
		limit = len(s.events)
	}
	out := make([]event.Event, 0, limit)
	for i := len(s.events) - 1; i >= 0 && len(out) < limit; i-- {
		out = append(out, s.events[i])
	}
	return out, nil
}
func (s *Store) Put(ctx context.Context, a *alert.Alert) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alerts[a.ID] = a
	s.byFingerprint[a.Fingerprint] = a.ID
	return nil
}
func (s *Store) Get(ctx context.Context, id string) (*alert.Alert, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.alerts[id]
	if !ok {
		return nil, fmt.Errorf("alert %s not found", id)
	}
	return a, nil
}
func (s *Store) FindByFingerprint(ctx context.Context, fp string) (*alert.Alert, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	id := s.byFingerprint[fp]
	if id == "" {
		return nil, fmt.Errorf("alert fingerprint not found")
	}
	return s.alerts[id], nil
}
func (s *Store) List(ctx context.Context, status alert.Status, limit int) ([]*alert.Alert, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]*alert.Alert, 0)
	for _, a := range s.alerts {
		if status != "" && a.Status != status {
			continue
		}
		out = append(out, a)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].LastSeen.After(out[j].LastSeen) })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}
func (s *Store) Save(ctx context.Context, r rule.Rule) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules[r.ID] = r
	return nil
}
func (s *Store) GetRule(ctx context.Context, id string) (rule.Rule, error) {
	if err := ctx.Err(); err != nil {
		return rule.Rule{}, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rules[id]
	if !ok {
		return r, fmt.Errorf("rule %s not found", id)
	}
	return r, nil
}
func (s *Store) ListRules(ctx context.Context, enabled bool) ([]rule.Rule, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []rule.Rule{}
	for _, r := range s.rules {
		if enabled && !r.Enabled {
			continue
		}
		out = append(out, r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(out[j].CreatedAt) })
	return out, nil
}
func (s *Store) Delete(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[id]; !ok {
		return fmt.Errorf("rule %s not found", id)
	}
	delete(s.rules, id)
	return nil
}
func (s *Store) SaveSilence(ctx context.Context, x silence.Silence) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.silences[x.ID] = x
	return nil
}
func (s *Store) ListSilences(ctx context.Context, active bool) ([]silence.Silence, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []silence.Silence{}
	for _, x := range s.silences {
		if active && !x.Active(timeNow()) {
			continue
		}
		out = append(out, x)
	}
	return out, nil
}
func timeNow() time.Time { return time.Now().UTC() }
func (s *Store) SaveNotification(ctx context.Context, x notification.Record) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notifications[x.ID] = x
	return nil
}
func (s *Store) UpdateNotification(ctx context.Context, x notification.Record) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notifications[x.ID] = x
	return nil
}
func (s *Store) ListNotifications(ctx context.Context, alertID string, limit int) ([]notification.Record, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := []notification.Record{}
	for _, x := range s.notifications {
		if alertID != "" && x.AlertID != alertID {
			continue
		}
		out = append(out, x)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].CreatedAt.Before(*out[j].CreatedAt) })
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}
