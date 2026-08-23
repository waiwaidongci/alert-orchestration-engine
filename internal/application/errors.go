package application

import "errors"

// Sentinel errors are the application-level vocabulary for failure modes that
// the HTTP layer must translate into status codes. Infrastructure adapters
// (e.g. the memory store) wrap these with %w so callers can match with
// errors.Is instead of grepping messages.
var (
	ErrNotFound     = errors.New("resource not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrConflict     = errors.New("conflict")
)
