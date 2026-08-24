package domain

import (
	"errors"
	"fmt"
	"strings"
)

// MaxNoteLength bounds how long an acknowledgement note may be. Notes are meant
// to be a short operator comment, so a generous cap keeps storage bounded while
// still allowing meaningful context.
const MaxNoteLength = 1000

type Gate struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Basin           string `json:"basin"`
	State           string `json:"state"`
	LastObservation string `json:"lastObservation"`
	AlertLevel      string `json:"alertLevel"`
	Acknowledged    bool   `json:"acknowledged"`
}

type AcknowledgeRequest struct {
	Note string `json:"note"`
}

func ValidateNote(note string) error {
	trimmed := strings.TrimSpace(note)
	if trimmed == "" {
		return errors.New("note is required")
	}
	if len(trimmed) > MaxNoteLength {
		return fmt.Errorf("note must be at most %d characters", MaxNoteLength)
	}
	return nil
}
