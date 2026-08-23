package httpadapter

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func decode(r *http.Request, v any) error {
	if r == nil {
		return fmt.Errorf("request is required")
	}
	body := r.Body
	if body == nil {
		body = io.NopCloser(strings.NewReader(""))
	}
	d := json.NewDecoder(io.LimitReader(body, 1<<20))
	d.DisallowUnknownFields()
	return d.Decode(v)
}
func errJSON(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, map[string]any{"error": map[string]any{"code": status, "message": err.Error()}})
}

func Decode(r *http.Request, v any) error { return decode(r, v) }
