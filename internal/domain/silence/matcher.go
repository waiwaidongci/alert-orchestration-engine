package silence

import (
	"sort"
	"strings"
	"time"
)

type Window struct{ Start, End time.Time }

func (w Window) Valid() bool { return !w.End.Before(w.Start) }
func (s Silence) Scope() string {
	parts := []string{}
	for k, v := range s.MatchLabels {
		parts = append(parts, k+"="+v)
	}
	sort.Strings(parts)
	return strings.Join(parts, ",")
}
func ActiveSilence(items []Silence, labels map[string]string, now time.Time) (Silence, bool) {
	for _, s := range items {
		if s.Active(now) && s.Matches(labels) {
			s.MatchLabels = cloneLabels(s.MatchLabels)
			return s, true
		}
	}
	return Silence{}, false
}

func cloneLabels(in map[string]string) map[string]string {
	out := make(map[string]string, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}
