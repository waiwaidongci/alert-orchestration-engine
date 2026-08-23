package application

import (
	"context"
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/domain/event"
	"github.com/example/alert-orchestration-engine/internal/domain/rule"
	"strings"
)

func ValidateEventInput(_ context.Context, in event.Input) error {
	if strings.TrimSpace(in.Source) == "" {
		return fmt.Errorf("source is required")
	}
	if strings.TrimSpace(in.Name) == "" {
		return fmt.Errorf("name is required")
	}
	if err := event.ValidateSource(in.Source); err != nil {
		return err
	}
	if err := event.ValidateName(in.Name); err != nil {
		return err
	}
	if in.Severity != "" {
		if err := event.ValidateSeverity(in.Severity); err != nil {
			return err
		}
	}
	return event.ValidateLabels(in.Labels)
}
func ValidateRuleInput(r rule.Rule) error {
	if err := r.Validate(); err != nil {
		return err
	}
	if err := rule.ValidateWindow(r.Window); err != nil {
		return err
	}
	r.Channels = rule.NormalizeChannels(r.Channels)
	return nil
}
func NormalizeLimit(limit int) int {
	if limit <= 0 {
		return 100
	}
	if limit > 500 {
		return 500
	}
	return limit
}
func NormalizeStatus(status string) string { return strings.ToLower(strings.TrimSpace(status)) }
