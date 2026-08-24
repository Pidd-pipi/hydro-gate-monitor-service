package store

import (
	"testing"

	"example.com/hydro-gate-monitor-service/domain"
)

func TestGateAckStoreMissingNoPanic(t *testing.T) {
	s := New()
	err := s.Acknowledge("gate-999", "note")
	if err == nil {
		t.Fatal("expected error for missing gate")
	}
	if err != domain.ErrGateNotFound {
		t.Fatalf("expected ErrGateNotFound, got %v", err)
	}
}

func TestGateAckListNoCrash(t *testing.T) {
	s := New()
	items := s.List()
	if len(items) == 0 {
		t.Fatal("list should still return seeded gates")
	}
}
