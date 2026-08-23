package httpadapter

import (
	"net/http"
)

type HealthState struct{ ready uint32 }

func (h *HealthState) SetReady(v bool) {
	if v {
		h.ready = 1
	} else {
		h.ready = 0
	}
}
func (h *HealthState) Ready() bool { return h.ready == 1 }
func (h *HealthState) Handler(w http.ResponseWriter, _ *http.Request) {
	if !h.Ready() {
		writeJSON(w, 503, map[string]string{"status": "starting"})
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ready"})
}
