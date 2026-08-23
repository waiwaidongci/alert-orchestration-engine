package alert

import (
	"sort"
	"strings"
	"time"
)

type Filter struct {
	Status     Status
	Source     string
	Severity   string
	LabelKey   string
	LabelValue string
	Since      *time.Time
	Until      *time.Time
}

func (a *Alert) Matches(f Filter) bool {
	if f.Status != "" && a.Status != f.Status {
		return false
	}
	if f.Source != "" && !strings.EqualFold(a.Source, f.Source) {
		return false
	}
	if f.Severity != "" && !strings.EqualFold(a.Severity, f.Severity) {
		return false
	}
	if f.LabelKey != "" && a.Labels[f.LabelKey] != f.LabelValue {
		return false
	}
	if f.Since != nil && a.LastSeen.Before(*f.Since) {
		return false
	}
	if f.Until != nil && a.LastSeen.After(*f.Until) {
		return false
	}
	return true
}
func SortByRecent(items []*Alert) []*Alert {
	out := append([]*Alert(nil), items...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].LastSeen.After(out[j].LastSeen) })
	return out
}
func CountByStatus(items []*Alert) map[Status]int {
	result := map[Status]int{}
	for _, a := range items {
		result[a.Status]++
	}
	return result
}
func IsTerminal(status Status) bool { return status == Resolved }
func AllowedTransition(from, to Status) bool {
	if from == Resolved {
		return true
	}
	switch to {
	case Acknowledged, Resolved, Suppressed:
		return true
	default:
		return false
	}
}
