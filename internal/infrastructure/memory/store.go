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
func (s *Store) SaveEvent(_ context.Context, e event.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events = append(s.events, e)
	return nil
}
func (s *Store) ListEvents(_ context.Context, limit int) ([]event.Event, error) {
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
func (s *Store) Put(_ context.Context, a *alert.Alert) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.alerts[a.ID] = a
	s.byFingerprint[a.Fingerprint] = a.ID
	return nil
}
func (s *Store) Get(_ context.Context, id string) (*alert.Alert, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	a, ok := s.alerts[id]
	if !ok {
		return nil, fmt.Errorf("alert %s not found", id)
	}
	return a, nil
}
func (s *Store) FindByFingerprint(_ context.Context, fp string) (*alert.Alert, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	id := s.byFingerprint[fp]
	if id == "" {
		return nil, fmt.Errorf("alert fingerprint not found")
	}
	return s.alerts[id], nil
}
func (s *Store) List(_ context.Context, status alert.Status, limit int) ([]*alert.Alert, error) {
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
func (s *Store) Save(_ context.Context, r rule.Rule) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.rules[r.ID] = r
	return nil
}
func (s *Store) GetRule(_ context.Context, id string) (rule.Rule, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	r, ok := s.rules[id]
	if !ok {
		return r, fmt.Errorf("rule %s not found", id)
	}
	return r, nil
}
func (s *Store) ListRules(_ context.Context, enabled bool) ([]rule.Rule, error) {
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
func (s *Store) Delete(_ context.Context, id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.rules[id]; !ok {
		return fmt.Errorf("rule %s not found", id)
	}
	delete(s.rules, id)
	return nil
}
func (s *Store) SaveSilence(_ context.Context, x silence.Silence) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.silences[x.ID] = x
	return nil
}
func (s *Store) ListSilences(_ context.Context, active bool) ([]silence.Silence, error) {
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
func (s *Store) SaveNotification(_ context.Context, x notification.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notifications[x.ID] = x
	return nil
}
func (s *Store) UpdateNotification(_ context.Context, x notification.Record) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.notifications[x.ID] = x
	return nil
}
func (s *Store) ListNotifications(_ context.Context, alertID string, limit int) ([]notification.Record, error) {
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
