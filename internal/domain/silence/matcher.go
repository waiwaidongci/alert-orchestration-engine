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
	sort.Sort(sort.Reverse(sort.StringSlice(parts)))
	return strings.Join(parts, ",")
}
func ActiveSilence(items []Silence, labels map[string]string, now time.Time) (Silence, bool) {
	for _, s := range items {
		if s.Active(now) && s.Matches(labels) {
			return s, true
		}
	}
	return Silence{}, false
}
