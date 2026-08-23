package memory

import (
	"context"
	"github.com/example/alert-orchestration-engine/internal/domain/alert"
	"github.com/example/alert-orchestration-engine/internal/domain/event"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
	"github.com/example/alert-orchestration-engine/internal/domain/rule"
	"github.com/example/alert-orchestration-engine/internal/domain/silence"
)

type EventRepo struct{ S *Store }

func (r EventRepo) SaveEvent(c context.Context, e event.Event) error { return r.S.SaveEvent(c, e) }
func (r EventRepo) ListEvents(c context.Context, n int) ([]event.Event, error) {
	return r.S.ListEvents(c, n)
}

type AlertRepo struct{ S *Store }

func (r AlertRepo) Put(c context.Context, a *alert.Alert) error            { return r.S.Put(c, a) }
func (r AlertRepo) Get(c context.Context, id string) (*alert.Alert, error) { return r.S.Get(c, id) }
func (r AlertRepo) FindByFingerprint(c context.Context, f string) (*alert.Alert, error) {
	return r.S.FindByFingerprint(c, f)
}
func (r AlertRepo) List(c context.Context, st alert.Status, n int) ([]*alert.Alert, error) {
	return r.S.List(c, st, n)
}

type RuleRepo struct{ S *Store }

func (r RuleRepo) Save(c context.Context, x rule.Rule) error           { return r.S.Save(c, x) }
func (r RuleRepo) Get(c context.Context, id string) (rule.Rule, error) { return r.S.GetRule(c, id) }
func (r RuleRepo) List(c context.Context, e bool) ([]rule.Rule, error) { return r.S.ListRules(c, e) }
func (r RuleRepo) Delete(c context.Context, id string) error           { return r.S.Delete(c, id) }

type SilenceRepo struct{ S *Store }

func (r SilenceRepo) Save(c context.Context, x silence.Silence) error { return r.S.SaveSilence(c, x) }
func (r SilenceRepo) List(c context.Context, a bool) ([]silence.Silence, error) {
	return r.S.ListSilences(c, a)
}

type NotificationRepo struct{ S *Store }

func (r NotificationRepo) Save(c context.Context, x notification.Record) error {
	return r.S.SaveNotification(c, x)
}
func (r NotificationRepo) Update(c context.Context, x notification.Record) error {
	return r.S.UpdateNotification(c, x)
}
func (r NotificationRepo) List(c context.Context, id string, n int) ([]notification.Record, error) {
	return r.S.ListNotifications(c, id, n)
}
