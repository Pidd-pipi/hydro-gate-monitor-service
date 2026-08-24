package health

import (
	"net/http"
	"sync/atomic"
)

var healthHits uint64

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeJSON(w, http.StatusMethodNotAllowed, `{"error":"method not allowed"}`)
		return
	}
	atomic.AddUint64(&healthHits, 1)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok","service":"hydro-gate-monitor-service"}`))
}

// Hits returns the total number of successful health checks observed so far.
func Hits() uint64 { return atomic.LoadUint64(&healthHits) }

func writeJSON(w http.ResponseWriter, status int, body string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write([]byte(body))
}
