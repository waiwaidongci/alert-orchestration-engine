package application

import (
	"context"
	"github.com/example/alert-orchestration-engine/internal/domain/alert"
	"github.com/example/alert-orchestration-engine/internal/domain/event"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
	"github.com/example/alert-orchestration-engine/internal/domain/rule"
	"github.com/example/alert-orchestration-engine/internal/domain/silence"
	"time"
)

type EventRepository interface {
	SaveEvent(context.Context, event.Event) error
	ListEvents(context.Context, int) ([]event.Event, error)
}
type AlertRepository interface {
	Put(context.Context, *alert.Alert) error
	Get(context.Context, string) (*alert.Alert, error)
	FindByFingerprint(context.Context, string) (*alert.Alert, error)
	List(context.Context, alert.Status, int) ([]*alert.Alert, error)
}
type RuleRepository interface {
	Save(context.Context, rule.Rule) error
	Get(context.Context, string) (rule.Rule, error)
	List(context.Context, bool) ([]rule.Rule, error)
	Delete(context.Context, string) error
}
type SilenceRepository interface {
	Save(context.Context, silence.Silence) error
	List(context.Context, bool) ([]silence.Silence, error)
}
type NotificationRepository interface {
	Save(context.Context, notification.Record) error
	Update(context.Context, notification.Record) error
	List(context.Context, string, int) ([]notification.Record, error)
}
type Notifier interface {
	Send(context.Context, notification.Record, *alert.Alert) error
}
type Clock interface{ Now() time.Time }

type IDGenerator interface{ NewID(string) string }
