package main

import "testing"

func TestRuleA0101LabelsDetached(t *testing.T) {
	first := opsRule0101()
	second := opsRule0103()
	first.RequiredLabels[0] = "a"
	second.RequiredLabels[0] = "b"
	again1 := opsRule0101()
	again3 := opsRule0103()
	if again1.RequiredLabels[0] != "site" || again3.RequiredLabels[0] != "site" {
		t.Fatalf("rules 0101/0103 share backing: %v %v", again1.RequiredLabels, again3.RequiredLabels)
	}
}

func TestRuleA0105LabelsDetached(t *testing.T) {
	first := opsRule0105()
	second := opsRule0107()
	first.RequiredLabels[0] = "a"
	second.RequiredLabels[0] = "b"
	again1 := opsRule0105()
	again7 := opsRule0107()
	if again1.RequiredLabels[0] != "site" || again7.RequiredLabels[0] != "site" {
		t.Fatalf("rules 0105/0107 share backing: %v %v", again1.RequiredLabels, again7.RequiredLabels)
	}
}

func TestRuleA0201LabelsDetached(t *testing.T) {
	first := opsRule0201()
	second := opsRule0203()
	first.RequiredLabels[0] = "a"
	second.RequiredLabels[0] = "b"
	again1 := opsRule0201()
	again3 := opsRule0203()
	if again1.RequiredLabels[0] != "site" || again3.RequiredLabels[0] != "site" {
		t.Fatalf("rules 0201/0203 share backing: %v %v", again1.RequiredLabels, again3.RequiredLabels)
	}
}

func TestRuleA0205LabelsDetached(t *testing.T) {
	first := opsRule0205()
	second := opsRule0207()
	first.RequiredLabels[0] = "a"
	second.RequiredLabels[0] = "b"
	again1 := opsRule0205()
	again7 := opsRule0207()
	if again1.RequiredLabels[0] != "site" || again7.RequiredLabels[0] != "site" {
		t.Fatalf("rules 0205/0207 share backing: %v %v", again1.RequiredLabels, again7.RequiredLabels)
	}
}
