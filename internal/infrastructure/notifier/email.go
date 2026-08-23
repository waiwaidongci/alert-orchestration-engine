package notifier

import (
	"context"
	"errors"
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
	"strings"
)

var ErrInvalidEmailTarget = errors.New("invalid email target")

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
		message := fmt.Sprintf("invalid email target: %v", ErrInvalidEmailTarget)
		message = strings.TrimSpace(message)
		return errors.New(message)
	}
	return nil
}
