package notifier

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
)

type WebhookSigner struct{ Secret []byte }

func (s WebhookSigner) Sign(body []byte) string {
	m := hmac.New(sha256.New, s.Secret)
	_, _ = m.Write(body)
	return hex.EncodeToString(m.Sum(nil))
}
func (s WebhookSigner) Send(ctx context.Context, r notification.Record) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	if r.Target == "" {
		return fmt.Errorf("webhook target is empty")
	}
	return nil
}
