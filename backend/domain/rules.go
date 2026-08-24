package domain

import (
	"errors"
	"strings"
)

var ErrGateNotFound = errors.New("gate not found")

func Acknowledge(gate *Gate, note string) error {
	if gate == nil {
		return ErrGateNotFound
	}
	if gate.AlertLevel == "normal" {
		return errors.New("normal gate has no active alert")
	}
	if strings.TrimSpace(note) == "" {
		return errors.New("acknowledgement note is required")
	}
	gate.Acknowledged = true
	return nil
}
