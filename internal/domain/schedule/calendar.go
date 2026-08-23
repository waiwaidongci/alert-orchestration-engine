package schedule

import (
	"fmt"
	"time"
)

type Rotation struct {
	Members  []string
	Interval time.Duration
	Anchor   time.Time
}

func (r Rotation) Member(now time.Time) string {
	if len(r.Members) == 0 {
		return ""
	}
	if r.Interval <= 0 {
		return r.Members[0]
	}
	steps := int(now.Sub(r.Anchor) / r.Interval)
	if steps < 0 {
		steps = 0
	}
	return r.Members[steps%len(r.Members)]
}
func (s DutySchedule) IsOnCall(now time.Time) bool { return s.Enabled && s.Covers(now.Hour()) }
func (s DutySchedule) Validate() error {
	if s.StartHour < 0 || s.StartHour > 23 || s.EndHour < 0 || s.EndHour > 23 {
		return fmt.Errorf("hours must be between 0 and 23")
	}
	return nil
}
