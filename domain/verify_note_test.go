package domain

import (
	"strings"
	"testing"
)

func TestConfigNoteValidateNoteRejectsBlank(t *testing.T) {
	if err := ValidateNote("   "); err == nil {
		t.Fatal("whitespace-only note should be rejected")
	}
}

func TestConfigNoteValidateNoteRejectsOversized(t *testing.T) {
	long := strings.Repeat("a", 100000)
	if err := ValidateNote(long); err == nil {
		t.Fatal("oversized note should be rejected")
	}
}
