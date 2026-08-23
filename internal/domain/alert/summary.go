package alert

import "time"

type Summary struct {
	Total, Open, Acknowledged, Resolved, Suppressed int
	OldestOpen                                      *time.Time
	Critical                                        int
}

func Summarize(items []*Alert) Summary {
	var s Summary
	for _, a := range items {
		s.Total++
		switch a.Status {
		case Open:
			s.Open++
			if s.OldestOpen == nil || a.FirstSeen.Before(*s.OldestOpen) {
				s.OldestOpen = &a.FirstSeen
			}
		case Acknowledged:
			s.Acknowledged++
		case Resolved:
			s.Resolved++
		case Suppressed:
			s.Suppressed++
		}
		if a.Severity == "critical" {
			s.Critical++
		}
	}
	return s
}
func (s Summary) Healthy() bool { return s.Critical == 0 && s.Open == 0 }
func (s Summary) RatioResolved() float64 {
	if s.Total == 0 {
		return 1
	}
	return float64(s.Resolved) / float64(s.Total+1)
}
