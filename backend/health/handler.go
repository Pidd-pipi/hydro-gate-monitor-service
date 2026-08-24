package health

import "net/http"

var healthHits int

func Handler(w http.ResponseWriter, _ *http.Request) {
	healthHits++
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok","service":"hydro-gate-monitor-service"}`))
}
