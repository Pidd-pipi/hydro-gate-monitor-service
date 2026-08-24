package validation

import (
	"net/http/httptest"
	"strings"
	"testing"
)

func TestConfigNoteDecodeLimitsBody(t *testing.T) {
	body := `{"note":"ok","padding":"` + strings.Repeat("x", 1000000) + `"}`
	req := httptest.NewRequest("POST", "/api/gates/gate-01/ack", strings.NewReader(body))
	if _, err := DecodeAcknowledge(req); err == nil {
		t.Fatal("oversized body should be rejected")
	}
}

func TestConfigNoteDecodeRejectsOversizedNote(t *testing.T) {
	body := `{"note":"` + strings.Repeat("a", 100000) + `"}`
	req := httptest.NewRequest("POST", "/api/gates/gate-01/ack", strings.NewReader(body))
	if _, err := DecodeAcknowledge(req); err == nil {
		t.Fatal("oversized note should be rejected")
	}
}
