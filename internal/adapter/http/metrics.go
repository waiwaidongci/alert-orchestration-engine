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

func (m *Metrics) Request() { atomic.AddUint64(&m.requests, 1) }
func (m *Metrics) Event()   { atomic.AddUint64(&m.events, 1) }
func (m *Metrics) Handler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")
	fmt.Fprintf(w, "alert_http_requests_total %d\nalert_events_ingested_total %d\n", atomic.LoadUint64(&m.requests), atomic.LoadUint64(&m.events))
}
