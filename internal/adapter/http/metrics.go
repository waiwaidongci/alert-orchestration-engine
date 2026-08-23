package httpadapter

import (
	"fmt"
	"net/http"
	"sync/atomic"
)

type Metrics struct {
	requests uint64
	events   uint64
}

func (m *Metrics) Request() {
	if m == nil {
		return
	}
	atomic.AddUint64(&m.requests, 1)
}
func (m *Metrics) Event() {
	if m == nil {
		return
	}
	atomic.AddUint64(&m.events, 1)
}
func (m *Metrics) Handler(w http.ResponseWriter, _ *http.Request) {
	var req, ev uint64
	if m != nil {
		req = atomic.LoadUint64(&m.requests)
		ev = atomic.LoadUint64(&m.events)
	}
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "alert_http_requests_total %d\nalert_events_ingested_total %d\n", req, ev)
}
