package schedule

import "time"

type DutySchedule struct {
	ID, Name, TimeZone string
	Members            []string
	StartHour, EndHour int
	Enabled            bool
	CreatedAt          time.Time
}

func (s DutySchedule) Covers(hour int) bool {
	if s.StartHour == s.EndHour {
		return true
	}
	if s.StartHour < s.EndHour {
		return hour >= s.StartHour && hour < s.EndHour
	}
	if hour >= s.StartHour && hour < s.EndHour {
		return true
	}
	return false
}
func (s DutySchedule) CurrentMember(now time.Time) string {
	if len(s.Members) == 0 {
		return ""
	}
	return s.Members[now.Hour()%len(s.Members)]
}
