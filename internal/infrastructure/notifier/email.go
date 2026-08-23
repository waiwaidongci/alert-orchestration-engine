package notifier

import (
	"context"
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
	"strings"
)

type EmailSender struct{ From string }

func (e EmailSender) Send(ctx context.Context, r notification.Record) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if e.From == "" {
		return fmt.Errorf("sender address is empty")
	}
	if !strings.Contains(r.Target, "@") {
		return fmt.Errorf("invalid email target")
	}
	return nil
}
