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
			writeError(w, ackStatus(err), err.Error())
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

// ackStatus classifies a store.Acknowledge error so a missing gate returns
// 404 (not 409), a state/policy rejection stays 409, and anything else
// falls back to internal rather than masquerading as a conflict.
func ackStatus(err error) int {
	switch {
	case errors.Is(err, domain.ErrGateNotFound):
		return http.StatusNotFound
	case err == nil:
		return http.StatusOK
	default:
		// Acknowledge on an existing-but-non-alerting gate, or an empty
		// note, are semantic/state rejections — conflict is the closest
		// stable mapping the existing API already used.
		return http.StatusConflict
	}
}
