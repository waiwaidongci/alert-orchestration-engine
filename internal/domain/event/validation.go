package event

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"
	"time"
)

var namePattern = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_.:-]{1,127}$`)

func ValidateSource(source string) error {
	source = strings.TrimSpace(source)
	if source == "" {
		return fmt.Errorf("source must not be empty")
	}
	if len(source) > 128 {
		return fmt.Errorf("source exceeds 128 characters")
	}
	if strings.ContainsAny(source, "\r\n") {
		return fmt.Errorf("source contains a newline")
	}
	return nil
}

func ValidateName(name string) error {
	if !namePattern.MatchString(strings.TrimSpace(name)) {
		return fmt.Errorf("name must match %s", namePattern.String())
	}
	return nil
}

func ValidateSeverity(severity string) error {
	switch strings.ToLower(strings.TrimSpace(severity)) {
	case "info", "warning", "critical", "error":
		return nil
	default:
		return fmt.Errorf("unsupported severity %q", severity)
	}
}

func ValidateLabels(labels map[string]string) error {
	if len(labels) > 32 {
		return fmt.Errorf("at most 32 labels are allowed")
	}
	for key, value := range labels {
		if strings.TrimSpace(key) == "" {
			return fmt.Errorf("label key is empty")
		}
		if len(key) > 64 || len(value) > 256 {
			return fmt.Errorf("label %s exceeds size limit", key)
		}
	}
	return nil
}

func ParseOccurredAt(value string) (time.Time, error) {
	if value == "" {
		return time.Time{}, nil
	}
	t, err := time.Parse(time.RFC3339Nano, value)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse occurred_at: %v", err)
	}
	return t, nil
}
func ValidateCallback(raw string) error {
	if raw == "" {
		return nil
	}
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "https" {
		return fmt.Errorf("callback must be an https URL")
	}
	return nil
}
