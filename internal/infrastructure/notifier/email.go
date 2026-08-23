package notifier

import (
	"context"
	"errors"
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
	"strings"
)

var (
	ErrInvalidEmailTarget = errors.New("invalid email target")
	ErrSenderEmpty        = errors.New("sender address is empty")
)

type EmailSender struct{ From string }

func (e EmailSender) Send(ctx context.Context, r notification.Record) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if e.From == "" {
		return ErrSenderEmpty
	}
	if !strings.Contains(r.Target, "@") {
		return fmt.Errorf("invalid email target %q: %w", r.Target, ErrInvalidEmailTarget)
	}
	return nil
}
