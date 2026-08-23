package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"example.com/hydro-gate-monitor-service/domain"
	"example.com/hydro-gate-monitor-service/store"
	"example.com/hydro-gate-monitor-service/validation"
)

func listGates(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			writeError(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		writeJSON(w, http.StatusOK, map[string]any{"gates": s.List()})
	}
}

func changeGate(s *store.Store) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || !strings.HasSuffix(r.URL.Path, "/ack") {
			writeError(w, http.StatusNotFound, "endpoint not found")
			return
		}
		input, err := validation.DecodeAcknowledge(r)
		if err != nil {
			writeError(w, http.StatusBadRequest, err.Error())
			return
		}
		if err := s.Acknowledge(pathID(r.URL.Path), input.Note); err != nil {
			status := http.StatusConflict
			if errors.Is(err, domain.ErrGateNotFound) {
				status = http.StatusNotFound
			}
			writeError(w, status, err.Error())
			return
		}
		writeJSON(w, http.StatusOK, map[string]string{"status": "acknowledged", "gateID": pathID(r.URL.Path)})
	}
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}
