package r009

import (
	"github.com/example/alert-orchestration-engine/internal/domain/event"
	"github.com/example/alert-orchestration-engine/internal/domain/rule"
	"github.com/example/alert-orchestration-engine/internal/domain/schedule"
	"testing"
	"time"
)

func TestR009Any(t *testing.T) {
	if (rule.Expression{Any: []rule.Condition{{Field: "severity", Operator: "eq", Value: "critical"}}}).Evaluate(event.Event{Severity: "warning"}) {
		t.Fatal("Any false")
	}
}
func TestR009Threshold(t *testing.T) {
	n := 10.
	if !(&rule.Rule{Enabled: true, EventName: "cpu", Threshold: &n}).Matches(event.Event{Name: "cpu", Value: 10}) {
		t.Fatal("equal threshold rejected")
	}
}
func TestR009Rotation(t *testing.T) {
	r := schedule.Rotation{Members: []string{"a", "b"}, Interval: time.Hour, Anchor: time.Unix(3600, 0)}
	if r.Member(time.Unix(0, 0)) != "a" {
		t.Fatal("pre-anchor advanced")
	}
}
func TestR009Schedule(t *testing.T) {
	s := schedule.DutySchedule{StartHour: 22, EndHour: 2}
	if !s.Covers(23) || !s.Covers(1) || s.Covers(12) {
		t.Fatal("cross midnight")
	}
}
