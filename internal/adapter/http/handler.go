package httpadapter

import (
	"fmt"
	"github.com/example/alert-orchestration-engine/internal/application"
	"github.com/example/alert-orchestration-engine/internal/domain/event"
	"github.com/example/alert-orchestration-engine/internal/domain/rule"
	"github.com/example/alert-orchestration-engine/internal/domain/silence"
	"net/http"
	"strconv"
)

type Handler struct {
	svc     *application.Service
	metrics *Metrics
	ready   func() bool
}

func NewHandler(svc *application.Service, m *Metrics, ready func() bool) *Handler {
	return &Handler{svc: svc, metrics: m, ready: ready}
}
func (h *Handler) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", h.health)
	mux.HandleFunc("GET /readyz", h.readyz)
	mux.HandleFunc("GET /metrics", h.metric)
	mux.HandleFunc("POST /api/v1/events", h.ingest)
	mux.HandleFunc("GET /api/v1/alerts", h.alerts)
	mux.HandleFunc("POST /api/v1/alerts/{id}/acknowledge", h.acknowledge)
	mux.HandleFunc("POST /api/v1/alerts/{id}/resolve", h.resolve)
	mux.HandleFunc("GET /api/v1/notifications", h.notifications)
	mux.HandleFunc("GET /api/v1/rules", h.rules)
	mux.HandleFunc("POST /api/v1/rules", h.createRule)
	mux.HandleFunc("DELETE /api/v1/rules/{id}", h.deleteRule)
	mux.HandleFunc("POST /api/v1/silences", h.createSilence)
	mux.HandleFunc("GET /api/v1/silences", h.silences)
	return mux
}
func (h *Handler) health(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, 200, map[string]string{"status": "ok"})
}
func (h *Handler) readyz(w http.ResponseWriter, _ *http.Request) {
	if h.ready != nil && !h.ready() {
		writeJSON(w, 503, map[string]string{"status": "not_ready"})
		return
	}
	writeJSON(w, 200, map[string]string{"status": "ready"})
}
func (h *Handler) metric(w http.ResponseWriter, r *http.Request) { h.metrics.Handler(w, r) }
func (h *Handler) ingest(w http.ResponseWriter, r *http.Request) {
	h.metrics.Request()
	var in event.Input
	if err := decode(r, &in); err != nil {
		errJSON(w, 400, fmt.Errorf("invalid event: %w", err))
		return
	}
	_, as, err := h.svc.Ingest(r.Context(), in)
	if err != nil {
		errJSON(w, 422, err)
		return
	}
	h.metrics.Event()
	writeJSON(w, 202, map[string]any{"accepted": true, "alerts": as})
}
func (h *Handler) alerts(w http.ResponseWriter, r *http.Request) {
	h.metrics.Request()
	as, err := h.svc.ListAlerts(r.Context(), r.URL.Query().Get("status"), parseLimit(r))
	if err != nil {
		errJSON(w, 500, err)
		return
	}
	writeJSON(w, 200, map[string]any{"items": as, "count": len(as)})
}
func (h *Handler) acknowledge(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := h.svc.Acknowledge(r.Context(), id)
	if err != nil {
		errJSON(w, 409, err)
		return
	}
	writeJSON(w, 200, a)
}
func (h *Handler) resolve(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	a, err := h.svc.Resolve(r.Context(), id)
	if err != nil {
		errJSON(w, 409, err)
		return
	}
	writeJSON(w, 200, a)
}
func (h *Handler) notifications(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.Notifications(r.Context(), r.URL.Query().Get("alert_id"), parseLimit(r))
	if err != nil {
		errJSON(w, 500, err)
		return
	}
	writeJSON(w, 200, map[string]any{"items": items, "count": len(items)})
}
func (h *Handler) rules(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.Rules(r.Context())
	if err != nil {
		errJSON(w, 500, err)
		return
	}
	writeJSON(w, 200, map[string]any{"items": items, "count": len(items)})
}
func (h *Handler) createRule(w http.ResponseWriter, r *http.Request) {
	var in rule.Rule
	if err := decode(r, &in); err != nil {
		errJSON(w, 400, err)
		return
	}
	x, err := h.svc.CreateRule(r.Context(), in)
	if err != nil {
		errJSON(w, 422, err)
		return
	}
	writeJSON(w, 201, x)
}
func (h *Handler) deleteRule(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	if err := h.svc.DeleteRule(r.Context(), id); err != nil {
		errJSON(w, 404, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
func (h *Handler) createSilence(w http.ResponseWriter, r *http.Request) {
	var in silence.Silence
	if err := decode(r, &in); err != nil {
		errJSON(w, 400, err)
		return
	}
	x, err := h.svc.CreateSilence(r.Context(), in)
	if err != nil {
		errJSON(w, 422, err)
		return
	}
	writeJSON(w, 201, x)
}
func (h *Handler) silences(w http.ResponseWriter, r *http.Request) {
	items, err := h.svc.Silences(r.Context())
	if err != nil {
		errJSON(w, 500, err)
		return
	}
	writeJSON(w, 200, map[string]any{"items": items, "count": len(items)})
}
func parseLimit(r *http.Request) int {
	n, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if n == 0 {
		n = 100
	}
	return n
}
