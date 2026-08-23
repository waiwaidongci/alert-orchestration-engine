package rule

import (
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/domain/event"
	"strings"
	"time"
)

type Condition struct {
	Field    string
	Operator string
	Value    string
}

func (c Condition) Evaluate(e event.Event) bool {
	var actual string
	switch c.Field {
	case "source":
		actual = e.Source
	case "name":
		actual = e.Name
	case "severity":
		actual = e.Severity
	default:
		actual = e.Labels[c.Field]
	}
	switch c.Operator {
	case "eq", "=":
		return actual == c.Value
	case "neq", "!=":
		return actual != c.Value
	case "contains":
		return strings.Contains(actual, c.Value)
	case "prefix":
		return strings.HasPrefix(actual, c.Value)
	default:
		return false
	}
}

type Expression struct {
	All []Condition
	Any []Condition
}

func (x Expression) Evaluate(e event.Event) bool {
	for _, c := range x.All {
		if !c.Evaluate(e) {
			return false
		}
	}
	if len(x.Any) == 0 {
		return true
	}
	for _, c := range x.Any {
		if c.Evaluate(e) {
			return true
		}
	}
	return false
}
func ValidateWindow(d time.Duration) error {
	if d < 0 {
		return fmt.Errorf("window must be non-negative")
	}
	if d > 7*24*time.Hour {
		return fmt.Errorf("window cannot exceed seven days")
	}
	return nil
}
func NormalizeChannels(channels []string) []string {
	seen := map[string]bool{}
	out := []string{}
	for _, c := range channels {
		c = strings.ToLower(strings.TrimSpace(c))
		if c != "" && !seen[c] {
			seen[c] = true
			out = append(out, c)
		}
	}
	if len(out) == 0 {
		return []string{"webhook"}
	}
	return out
}
