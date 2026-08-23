package application

import (
	"context"
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/domain/alert"
	"github.com/example/alert-orchestration-engine/internal/domain/event"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
	"github.com/example/alert-orchestration-engine/internal/domain/rule"
	"github.com/example/alert-orchestration-engine/internal/domain/silence"
	"strings"
	"time"
)

type Service struct {
	events        EventRepository
	alerts        AlertRepository
	rules         RuleRepository
	silences      SilenceRepository
	notifications NotificationRepository
	notifier      Notifier
	clock         Clock
	ids           IDGenerator
}

func NewService(e EventRepository, a AlertRepository, r RuleRepository, s SilenceRepository, n NotificationRepository, sender Notifier, c Clock, ids IDGenerator) *Service {
	return &Service{events: e, alerts: a, rules: r, silences: s, notifications: n, notifier: sender, clock: c, ids: ids}
}
func (s *Service) Ingest(ctx context.Context, in event.Input) (event.Event, []*alert.Alert, error) {
	now := s.clock.Now()
	e, err := event.Normalize(ctx, in, now)
	if err != nil {
		return event.Event{}, nil, err
	}
	if err := s.events.SaveEvent(ctx, e); err != nil {
		return e, nil, fmt.Errorf("save event: %w", err)
	}
	rs, err := s.rules.List(ctx, true)
	if err != nil {
		return e, nil, fmt.Errorf("list rules: %w", err)
	}
	created := make([]*alert.Alert, 0)
	for _, r := range rs {
		if !r.Matches(e) {
			continue
		}
		existing, _ := s.alerts.FindByFingerprint(ctx, e.Fingerprint())
		if existing != nil {
			existing.Touch(now, e)
			if err := s.alerts.Put(ctx, existing); err != nil {
				return e, nil, fmt.Errorf("update alert: %w", err)
			}
			created = append(created, existing)
			continue
		}
		a := alert.New(s.ids.NewID("alt"), e, r.ID, now)
		silenced, _ := s.isSilenced(ctx, e.Labels, now)
		if silenced {
			a.Status = alert.Suppressed
		}
		if err := s.alerts.Put(ctx, a); err != nil {
			return e, nil, fmt.Errorf("save alert: %w", err)
		}
		created = append(created, a)
		if !silenced {
			s.dispatch(ctx, a, r)
		}
	}
	return e, created, nil
}
func (s *Service) isSilenced(ctx context.Context, labels map[string]string, now time.Time) (bool, error) {
	ss, err := s.silences.List(ctx, true)
	if err != nil {
		return false, err
	}
	for _, x := range ss {
		if x.Active(now) && x.Matches(labels) {
			return true, nil
		}
	}
	return false, nil
}
func (s *Service) dispatch(ctx context.Context, a *alert.Alert, r rule.Rule) {
	for _, ch := range r.Channels {
		if err := ctx.Err(); err != nil {
			return
		}
		rec := notification.New(s.ids.NewID("ntf"), a.ID, r.ID, ch, ch, s.clock.Now())
		if err := s.notifications.Save(ctx, rec); err != nil {
			continue
		}
		go s.send(ctx, rec, a)
	}
}
func (s *Service) send(ctx context.Context, rec notification.Record, a *alert.Alert) {
	if err := ctx.Err(); err != nil {
		return
	}
	err := s.notifier.Send(ctx, rec, a)
	if err == nil {
		now := s.clock.Now()
		rec.Status = notification.Sent
		rec.Attempts++
		rec.SentAt = &now
	} else {
		rec.Status = notification.Failed
		rec.Attempts++
		rec.Error = err.Error()
	}
	_ = s.notifications.Update(context.Background(), rec)
}
func (s *Service) ListAlerts(ctx context.Context, status string, limit int) ([]*alert.Alert, error) {
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	return s.alerts.List(ctx, alert.Status(strings.ToLower(status)), limit)
}
func (s *Service) Acknowledge(ctx context.Context, id string) (*alert.Alert, error) {
	a, err := s.alerts.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := a.Transition(alert.Acknowledged, s.clock.Now()); err != nil {
		return nil, err
	}
	return a, s.alerts.Put(ctx, a)
}
func (s *Service) Resolve(ctx context.Context, id string) (*alert.Alert, error) {
	a, err := s.alerts.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := a.Transition(alert.Resolved, s.clock.Now()); err != nil {
		return nil, err
	}
	return a, s.alerts.Put(ctx, a)
}
func (s *Service) CreateRule(ctx context.Context, r rule.Rule) (rule.Rule, error) {
	if r.ID == "" {
		r.ID = s.ids.NewID("rul")
	}
	if r.CreatedAt.IsZero() {
		r.CreatedAt = s.clock.Now()
	}
	r.UpdatedAt = s.clock.Now()
	if r.Channels == nil {
		r.Channels = []string{"webhook"}
	}
	if err := r.Validate(); err != nil {
		return r, err
	}
	return r, s.rules.Save(ctx, r)
}
func (s *Service) Rules(ctx context.Context) ([]rule.Rule, error)  { return s.rules.List(ctx, false) }
func (s *Service) DeleteRule(ctx context.Context, id string) error { return s.rules.Delete(ctx, id) }
func (s *Service) CreateSilence(ctx context.Context, x silence.Silence) (silence.Silence, error) {
	if x.ID == "" {
		x.ID = s.ids.NewID("sil")
	}
	if x.CreatedAt.IsZero() {
		x.CreatedAt = s.clock.Now()
	}
	if x.EndsAt.IsZero() {
		x.EndsAt = x.CreatedAt.Add(time.Hour)
	}
	if x.MatchLabels == nil {
		x.MatchLabels = map[string]string{}
	}
	return x, s.silences.Save(ctx, x)
}
func (s *Service) Silences(ctx context.Context) ([]silence.Silence, error) {
	return s.silences.List(ctx, false)
}
func (s *Service) Notifications(ctx context.Context, alertID string, limit int) ([]notification.Record, error) {
	return s.notifications.List(ctx, alertID, limit)
}
