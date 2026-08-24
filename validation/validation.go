package validation

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"example.com/hydro-gate-monitor-service/domain"
)

func DecodeAcknowledge(r *http.Request) (domain.AcknowledgeRequest, error) {
	var input domain.AcknowledgeRequest
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&input); err != nil {
		return input, errors.New("body must be valid JSON")
	}
	if strings.TrimSpace(input.Note) == "" {
		return input, errors.New("note must not be empty")
	}
	return input, nil
}
