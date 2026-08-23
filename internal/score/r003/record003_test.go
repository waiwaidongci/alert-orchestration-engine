package r003

import (
	"github.com/example/alert-orchestration-engine/internal/domain/alert"
	"github.com/example/alert-orchestration-engine/internal/domain/event"
	"github.com/example/alert-orchestration-engine/internal/domain/notification"
	"github.com/example/alert-orchestration-engine/internal/domain/silence"
	"testing"
	"time"
)

func TestR003Transition(t *testing.T) {
	a := alert.New("a", event.Event{Source: "s", Name: "n"}, "r", time.Now())
	a.Status = alert.Resolved
	if a.Transition(alert.Acknowledged, time.Now()) == nil {
		t.Fatal("terminal transition accepted")
	}
}
func TestR003Table(t *testing.T) {
	if alert.AllowedTransition(alert.Resolved, alert.Acknowledged) {
		t.Fatal("table allowed rollback")
	}
}
func TestR003Sent(t *testing.T) {
	r := notification.New("n", "a", "r", "webhook", "x", time.Now())
	n := time.Now().Add(time.Hour)
	r.NextAttemptAt = &n
	r.MarkSent(time.Now())
	if r.NextAttemptAt != nil {
		t.Fatal("retry deadline retained")
	}
}
func TestR003Silence(t *testing.T) {
	e := time.Unix(100, 0)
	if (silence.Silence{StartsAt: e.Add(-time.Minute), EndsAt: e}).Active(e) {
		t.Fatal("end inclusive")
	}
}
