package validation

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"example.com/hydro-gate-monitor-service/domain"
)

// MaxRequestBodyBytes bounds the size of an acknowledgement request body so an
// oversized payload is rejected before it is fully buffered or decoded.
const MaxRequestBodyBytes = 1 << 14 // 16 KiB

func DecodeAcknowledge(r *http.Request) (domain.AcknowledgeRequest, error) {
	r.Body = http.MaxBytesReader(nil, r.Body, MaxRequestBodyBytes)
	var input domain.AcknowledgeRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		return input, errors.New("body must be valid JSON")
	}
	note := strings.TrimSpace(input.Note)
	if note == "" {
		return input, errors.New("note must not be empty")
	}
	if len(note) > domain.MaxNoteLength {
		return input, fmt.Errorf("note must be at most %d characters", domain.MaxNoteLength)
	}
	return input, nil
}
