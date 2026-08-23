package notification

import (
	"fmt"
	"math"
	"time"
)

const MaxAttempts = 5

func (r Record) CanRetry() bool { return r.Status == Failed && r.Attempts < MaxAttempts }
func RetryDelay(attempt int) time.Duration {
	if attempt < 1 {
		attempt = 1
	}
	if attempt > 10 {
		attempt = 10
	}
	return time.Duration(math.Pow(2, float64(attempt-1))) * time.Second
}
func (r *Record) MarkFailed(err error, now time.Time) {
	r.Status = Failed
	r.Attempts++
	if err != nil {
		r.Error = err.Error()
	}
	if r.Attempts >= MaxAttempts {
		r.Status = DeadLetter
		r.NextAttemptAt = nil
		return
	}
	next := now.Add(RetryDelay(r.Attempts))
	r.NextAttemptAt = &next
}
func (r *Record) MarkSent(now time.Time) {
	r.Status = Sent
	r.Attempts++
	r.Error = ""
	r.NextAttemptAt = nil
	r.SentAt = &now
}
func ValidateTarget(channel, target string) error {
	if target == "" {
		return fmt.Errorf("target is required for %s", channel)
	}
	return nil
}
