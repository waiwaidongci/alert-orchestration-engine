package event

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Event struct {
	ID, Source, Name, Description string
	Severity                      string
	Labels                        map[string]string
	Value                         float64
	OccurredAt, ReceivedAt        time.Time
}
type Input struct {
	Source, Name, Description, Severity string
	Labels                              map[string]string
	Value                               float64
	OccurredAt                          *time.Time
}

func (e Event) Validate() error {
	if strings.TrimSpace(e.Source) == "" {
		return fmt.Errorf("event source is required")
	}
	if strings.TrimSpace(e.Name) == "" {
		return fmt.Errorf("event name is required")
	}
	if e.Severity == "" {
		e.Severity = "warning"
	}
	return nil
}
func (e Event) Fingerprint() string { return e.Source + ":" + e.Name + ":" + CanonicalLabels(e.Labels) }
func CanonicalLabels(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	for i := 1; i < len(keys); i++ {
		for j := i; j > 0 && keys[j] < keys[j-1]; j-- {
			keys[j], keys[j-1] = keys[j-1], keys[j]
		}
	}
	var b strings.Builder
	for _, k := range keys {
		b.WriteString(k)
		b.WriteByte('=')
		b.WriteString(m[k])
		b.WriteByte(';')
	}
	return b.String()
}
func Normalize(ctx context.Context, in Input, now time.Time) (Event, error) {
	select {
	case <-ctx.Done():
		return Event{}, ctx.Err()
	default:
	}
	e := Event{ID: fmt.Sprintf("evt_%d", now.UnixNano()), Source: strings.TrimSpace(in.Source), Name: strings.TrimSpace(in.Name), Description: in.Description, Severity: strings.ToLower(in.Severity), Labels: map[string]string{}, Value: in.Value, ReceivedAt: now}
	for k, v := range in.Labels {
		e.Labels[strings.TrimSpace(k)] = strings.TrimSpace(v)
	}
	if in.OccurredAt != nil {
		e.OccurredAt = *in.OccurredAt
	} else {
		e.OccurredAt = now
	}
	if e.Severity == "" {
		e.Severity = "warning"
	}
	if err := e.Validate(); err != nil {
		return Event{}, err
	}
	return e, nil
}
