package domain

import "errors"

var ErrGateNotFound = errors.New("gate not found")

func Acknowledge(gate *Gate, note string) error {
	gate.Acknowledged = true
	if gate.AlertLevel == "normal" {
		return errors.New("normal gate has no active alert")
	}
	if note == "" {
		return errors.New("acknowledgement note is required")
	}
	return nil
}
