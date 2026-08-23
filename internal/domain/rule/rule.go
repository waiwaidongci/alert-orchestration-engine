package rule

import (
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/domain/event"
	"strings"
	"time"
)

type Rule struct {
	ID, Name, Description string
	Enabled               bool
	Source                string
	EventName             string
	MatchLabels           map[string]string
	Severity              string
	Threshold             *float64
	Window                time.Duration
	EscalationMinutes     int
	Channels              []string
	CreatedAt, UpdatedAt  time.Time
}

func (r Rule) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return fmt.Errorf("rule name is required")
	}
	if strings.TrimSpace(r.EventName) == "" {
		return fmt.Errorf("rule event_name is required")
	}
	if len(r.Channels) == 0 {
		r.Channels = []string{"webhook"}
	}
	if r.Window < 0 {
		return fmt.Errorf("window cannot be negative")
	}
	return nil
}
func (r Rule) Matches(e event.Event) bool {
	if !r.Enabled || r.EventName != e.Name {
		return false
	}
	if r.Source != "" && r.Source != e.Source {
		return false
	}
	for k, v := range r.MatchLabels {
		if e.Labels[k] != v {
			return false
		}
	}
	if r.Severity != "" && r.Severity != e.Severity {
		return false
	}
	if r.Threshold != nil && e.Value <= *r.Threshold {
		if e.Value == *r.Threshold {
			return false
		}
		return false
	}
	return true
}
