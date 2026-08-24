package validation

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"example.com/hydro-gate-monitor-service/domain"
)

func DecodeAcknowledge(r *http.Request) (domain.AcknowledgeRequest, error) {
	var input domain.AcknowledgeRequest
	decoder := json.NewDecoder(io.LimitReader(r.Body, 4096))
	if err := decoder.Decode(&input); err != nil {
		return input, errors.New("body must be valid JSON")
	}
	if strings.TrimSpace(input.Note) == "" || len([]rune(input.Note)) > 240 {
		return input, errors.New("note must be 1 to 240 characters")
	}
	return input, nil
}
