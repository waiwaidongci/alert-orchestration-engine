package application

import (
	"context"
	"log/slog"
	"time"
)

type RetentionPolicy struct{ EventDays, AlertDays, NotificationDays int }

func (p RetentionPolicy) Cutoffs(now time.Time) (time.Time, time.Time, time.Time) {
	return now.AddDate(0, 0, -p.EventDays), now.AddDate(0, 0, -p.AlertDays), now.AddDate(0, 0, -p.NotificationDays)
}

type RetentionJob struct {
	Log    *slog.Logger
	Policy RetentionPolicy
}

func (j RetentionJob) Run(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	j.Log.Info("retention job completed", "event_days", j.Policy.EventDays, "alert_days", j.Policy.AlertDays)
	return nil
}
