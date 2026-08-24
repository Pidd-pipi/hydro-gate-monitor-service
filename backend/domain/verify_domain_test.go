package domain

import "testing"

func TestGateAckInvalidNoteNoMutation(t *testing.T) {
	gate := &Gate{ID: "g1", AlertLevel: "watch"}
	if err := Acknowledge(gate, ""); err == nil {
		t.Fatal("empty note should be rejected")
	}
	if gate.Acknowledged {
		t.Fatal("rejected acknowledgement should not mutate gate state")
	}
}

func TestGateAckWhitespaceNoteRejected(t *testing.T) {
	gate := &Gate{ID: "g1", AlertLevel: "watch"}
	if err := Acknowledge(gate, "   "); err == nil {
		t.Fatal("whitespace-only note should be rejected")
	}
}
