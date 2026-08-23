package notification

import "time"

type Status string

const (
	Pending    Status = "pending"
	Sent       Status = "sent"
	Failed     Status = "failed"
	DeadLetter Status = "dead_letter"
)

type Record struct {
	ID, AlertID, RuleID, Channel, Target, Error string
	Status                                      Status
	Attempts                                    int
	CreatedAt, SentAt                           *time.Time
	NextAttemptAt                               *time.Time
}

func New(id, alertID, ruleID, channel, target string, now time.Time) Record {
	return Record{ID: id, AlertID: alertID, RuleID: ruleID, Channel: channel, Target: target, Status: Pending, Attempts: 0, CreatedAt: &now}
}
