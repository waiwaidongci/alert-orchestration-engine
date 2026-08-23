package notifier

import (
	"context"
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/domain/alert"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
	"strings"
	"sync"
)

type Dispatcher struct {
	mu   sync.Mutex
	Sent []notification.Record
}

func New() *Dispatcher { return &Dispatcher{Sent: []notification.Record{}} }
func (d *Dispatcher) Send(ctx context.Context, r notification.Record, a *alert.Alert) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if strings.TrimSpace(r.Channel) == "" {
		return fmt.Errorf("notification channel is required")
	}
	if r.Channel != "webhook" && r.Channel != "email" && r.Channel != "sms" {
		return fmt.Errorf("unsupported notification channel %q", r.Channel)
	}
	d.Sent = append(d.Sent, r)
	return nil
}
