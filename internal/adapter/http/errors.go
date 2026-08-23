package httpadapter

import (
	"errors"
	"net/http"

	"github.com/example/alert-orchestration-engine/internal/application"
)

// Sentinel errors are re-exports of the application-level vocabulary so the
// HTTP layer can map them to status codes. Callers may also match on the
// application-level sentinels directly; both resolve through errors.Is.
var (
	ErrNotFound      = application.ErrNotFound
	ErrInvalidInput  = application.ErrInvalidInput
	ErrConflict      = application.ErrConflict
	ErrInternal      = errors.New("internal error")
	ErrUnprocessable = errors.New("unprocessable")
)

// statusFor maps a domain/infra error to an HTTP status. It unwraps with
// errors.Is so callers may wrap the sentinels above (e.g.
// fmt.Errorf("alert %s: %w", id, application.ErrNotFound)) and still get the
// right code. Unknown errors default to 422 to preserve prior behavior.
func statusFor(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrInvalidInput):
		return http.StatusBadRequest
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, ErrInternal):
		return http.StatusInternalServerError
	default:
		return http.StatusUnprocessableEntity
	}
}

// errCode is the stable, machine-readable token surfaced in the JSON body so
// clients can branch on a code instead of grepping the message text.
func errCode(err error) string {
	switch {
	case errors.Is(err, ErrNotFound):
		return "not_found"
	case errors.Is(err, ErrInvalidInput):
		return "invalid_input"
	case errors.Is(err, ErrConflict):
		return "conflict"
	case errors.Is(err, ErrInternal):
		return "internal"
	default:
		return "unprocessable"
	}
}

func writeError(w http.ResponseWriter, err error) { errJSON(w, statusFor(err), err) }

func StatusFor(err error) int { return statusFor(err) }
