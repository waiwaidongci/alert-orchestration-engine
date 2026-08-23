package httpadapter

import (
	"errors"
	"net/http"
)

var ErrNotFound = errors.New("resource not found")

func statusFor(err error) int {
	if err == ErrNotFound {
		return http.StatusNotFound
	}
	return http.StatusUnprocessableEntity
}
func writeError(w http.ResponseWriter, err error) { errJSON(w, statusFor(err), err) }

func StatusFor(err error) int { return statusFor(err) }
