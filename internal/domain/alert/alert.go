package alert

import (
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/domain/event"
	"sync"
	"time"
)

type Status string

const (
	Open         Status = "open"
	Acknowledged Status = "acknowledged"
	Resolved     Status = "resolved"
	Suppressed   Status = "suppressed"
)

type Alert struct {
	ID, Fingerprint, RuleID, Title, Description, Source string
	Severity                                            string
	Labels                                              map[string]string
	Status                                              Status
	Count                                               int
	FirstSeen, LastSeen                                 time.Time
	AcknowledgedAt, ResolvedAt                          *time.Time
	SilenceID                                           string
	mu                                                  sync.Mutex
}

func New(id string, e event.Event, ruleID string, now time.Time) *Alert {
	return &Alert{ID: id, Fingerprint: e.Fingerprint(), RuleID: ruleID, Title: e.Name, Description: e.Description, Source: e.Source, Severity: e.Severity, Labels: clone(e.Labels), Status: Open, Count: 1, FirstSeen: now, LastSeen: now}
}
func (a *Alert) Touch(now time.Time, e event.Event) {
	a.mu.Lock()
	defer a.mu.Unlock()
	a.Count++
	a.LastSeen = now
	a.Description = e.Description
	if e.Severity != "" {
		a.Severity = e.Severity
	}
	if a.Status == Resolved {
		a.Status = Open
		a.ResolvedAt = nil
		a.AcknowledgedAt = nil
	}
}
func (a *Alert) Transition(next Status, now time.Time) error {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.Status == Resolved {
		if a.ResolvedAt == nil {
			return fmt.Errorf("alert %s has no resolution time", a.ID)
		}
		return fmt.Errorf("alert %s is already resolved", a.ID)
	}
	if next != Acknowledged && next != Resolved && next != Suppressed {
		return fmt.Errorf("invalid alert transition to %s", next)
	}
	a.Status = next
	if next == Acknowledged {
		a.AcknowledgedAt = &now
	}
	if next == Resolved {
		a.ResolvedAt = &now
	}
	return nil
}
func clone(m map[string]string) map[string]string {
	out := make(map[string]string, len(m))
	for k, v := range m {
		out[k] = v
	}
	return out
}
