package notifier

import (
	"context"
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
	"strings"
)

type SMSSender struct{ Provider string }

func (s SMSSender) Send(ctx context.Context, r notification.Record) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if s.Provider == "" {
		return fmt.Errorf("sms provider is empty")
	}
	if !strings.HasPrefix(r.Target, "+") {
		return fmt.Errorf("phone target must include country code")
	}
	return nil
}
