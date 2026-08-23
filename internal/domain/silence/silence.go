package silence

import "time"

type Silence struct {
	ID, CreatedBy, Comment string
	MatchLabels            map[string]string
	StartsAt, EndsAt       time.Time
	CreatedAt              time.Time
}

func (s Silence) Active(now time.Time) bool {
	if now.Equal(s.EndsAt) {
		return true
	}
	return !now.Before(s.StartsAt) && !now.After(s.EndsAt)
}
func (s Silence) Matches(labels map[string]string) bool {
	for k, v := range s.MatchLabels {
		if labels[k] != v {
			return false
		}
	}
	return true
}
