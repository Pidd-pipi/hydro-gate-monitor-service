package main

import "testing"

func TestRuleBTotalCountIs112(t *testing.T) {
	if got := len(opsRules()); got != 112 {
		t.Fatalf("rules count = %d, want 112", got)
	}
}

func TestRuleBCodesUnique(t *testing.T) {
	rules := opsRules()
	seen := map[string]bool{}
	for _, r := range rules {
		if seen[r.Code] {
			t.Fatalf("duplicate rule code %q", r.Code)
		}
		seen[r.Code] = true
	}
}

func TestRuleBLookupMissingReturnsNotFound(t *testing.T) {
	_, ok := opsRuleByCode("OPS-9999")
	if ok {
		t.Fatal("missing rule lookup should report not found")
	}
}

func TestRuleBLookupMissingDoesNotCrash(t *testing.T) {
	_, ok := opsFirstRequiredLabel("OPS-9999")
	if !ok {
		return
	}
}
