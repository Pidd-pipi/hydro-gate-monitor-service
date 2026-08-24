package main

import (
	"context"
	"testing"
)

func TestPageQueryDefaultsFillsZeroPage(t *testing.T) {
	q := opsQueryDefaults(OpsQuery{})
	if q.Page != 1 {
		t.Fatalf("page not defaulted: %d", q.Page)
	}
	if q.PageSize < 1 || q.PageSize > 200 {
		t.Fatalf("page size not defaulted: %d", q.PageSize)
	}
}

func TestPageSearchZeroPageNoPanic(t *testing.T) {
	svc := newOpsService([]OpsRecord{
		{ID: "r1", Subject: "subject a", Owner: "owner", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}},
		{ID: "r2", Subject: "subject b", Owner: "owner", Status: OpsStatusActive, Priority: OpsPriorityHigh, Labels: map[string]string{"site": "s1"}},
	})
	page, err := svc.Search(context.Background(), OpsQuery{PageSize: 25})
	if err != nil {
		t.Fatal(err)
	}
	if page.Page != 1 || len(page.Items) != 2 {
		t.Fatalf("zero page search failed: %+v", page)
	}
}

func TestPageBoundsClampsBeyondEnd(t *testing.T) {
	start, end := opsBounds(10, 99, 25)
	if start > 10 || start < 0 || end < start || end > 10 {
		t.Fatalf("bounds out of range: start=%d end=%d", start, end)
	}
	start2, end2 := opsBounds(0, 3, 25)
	if start2 != 0 || end2 != 0 {
		t.Fatalf("empty total should yield empty bounds: %d %d", start2, end2)
	}
}

func TestPageMoveFailureNoHistoryPollution(t *testing.T) {
	m := newOpsStateMachine()
	if err := m.Move(OpsStatusActive, OpsStatusQueued, "not allowed"); err == nil {
		t.Fatal("expected transition error")
	}
	if len(m.History()) != 0 {
		t.Fatalf("failed move polluted history: %+v", m.History())
	}
}

func TestPageHistorySnapshotIsolated(t *testing.T) {
	m := newOpsStateMachine()
	_ = m.Move(OpsStatusQueued, OpsStatusActive, "start")
	snapshot := m.History()
	snapshot[0].Reason = "tampered"
	latest, ok := m.Last()
	if !ok || latest.Reason == "tampered" {
		t.Fatalf("history snapshot aliases machine history: %+v", snapshot)
	}
}
