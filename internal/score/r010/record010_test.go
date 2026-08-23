package r010

import (
	"github.com/example/alert-orchestration-engine/internal/domain/alert"
	"github.com/example/alert-orchestration-engine/internal/domain/silence"
	"testing"
	"time"
)

func TestR010Scope(t *testing.T) {
	if (silence.Silence{MatchLabels: map[string]string{"b": "2", "a": "1"}}).Scope() != "a=1,b=2" {
		t.Fatal("scope order")
	}
}
func TestR010Copy(t *testing.T) {
	x := []silence.Silence{{MatchLabels: map[string]string{"env": "prod"}, StartsAt: time.Unix(0, 0), EndsAt: time.Unix(100, 0)}}
	got, ok := silence.ActiveSilence(x, map[string]string{"env": "prod"}, time.Unix(1, 0))
	if !ok {
		t.Fatal("inactive")
	}
	x[0].MatchLabels["env"] = "dev"
	if got.MatchLabels["env"] != "prod" {
		t.Fatal("map alias")
	}
}
func TestR010Summary(t *testing.T) {
	first := time.Unix(1, 0)
	a := &alert.Alert{Status: alert.Open, FirstSeen: first}
	s := alert.Summarize([]*alert.Alert{a})
	a.FirstSeen = time.Unix(9, 0)
	if s.OldestOpen == nil || !s.OldestOpen.Equal(first) {
		t.Fatal("pointer alias")
	}
}
func TestR010Ratio(t *testing.T) {
	if (alert.Summary{Total: 2, Resolved: 1}).RatioResolved() != 0.5 {
		t.Fatal("ratio")
	}
}
