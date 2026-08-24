package main

import (
	"context"
	"sync"
	"testing"
)

func TestStoreConcurrentListPutNoRace(t *testing.T) {
	s := newOpsStore(nil)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			rec := OpsRecord{ID: string(rune('a'+i)) + "1", Subject: "subject", Owner: "owner", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}}
			_ = s.Put(context.Background(), rec)
			_, _ = s.List(context.Background())
		}(i)
	}
	close(start)
	wg.Wait()
}

func TestStoreConcurrentGetUpdateNoRace(t *testing.T) {
	s := newOpsStore([]OpsRecord{{ID: "r1", Subject: "subject", Owner: "owner", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}}})
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			_, _ = s.Get(context.Background(), "r1")
			_ = s.Update(context.Background(), OpsRecord{ID: "r1", Subject: "subject2", Owner: "owner", Status: OpsStatusActive, Priority: OpsPriorityHigh, Labels: map[string]string{"site": "s2"}}, 0)
		}()
	}
	close(start)
	wg.Wait()
}

func TestStoreGetLabelsIndependent(t *testing.T) {
	s := newOpsStore([]OpsRecord{{ID: "r1", Subject: "subject", Owner: "owner", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}}})
	got, err := s.Get(context.Background(), "r1")
	if err != nil {
		t.Fatal(err)
	}
	got.Labels["site"] = "mutated"
	again, _ := s.Get(context.Background(), "r1")
	if again.Labels["site"] != "s1" {
		t.Fatalf("get labels not independent: %v", again.Labels)
	}
}

func TestStoreListReturnsIsolatedItems(t *testing.T) {
	s := newOpsStore([]OpsRecord{{ID: "r1", Subject: "subject", Owner: "owner", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}}})
	items, err := s.List(context.Background())
	if err != nil || len(items) != 1 {
		t.Fatalf("list failed: %v %v", items, err)
	}
	items[0].Labels["site"] = "mutated"
	again, _ := s.Get(context.Background(), "r1")
	if again.Labels["site"] != "s1" {
		t.Fatalf("list labels not independent: %v", again.Labels)
	}
}

func TestStorePutDoesNotAliasCallerLabels(t *testing.T) {
	s := newOpsStore(nil)
	labels := map[string]string{"site": "s1"}
	rec := OpsRecord{ID: "r1", Subject: "subject", Owner: "owner", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: labels}
	if err := s.Put(context.Background(), rec); err != nil {
		t.Fatal(err)
	}
	labels["site"] = "mutated"
	if got := s.items["r1"].Labels["site"]; got != "s1" {
		t.Fatalf("put aliased caller labels: %v", got)
	}
}

func TestStoreUpdateDoesNotAliasCallerLabels(t *testing.T) {
	s := newOpsStore([]OpsRecord{{ID: "r1", Subject: "subject", Owner: "owner", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}}})
	labels := map[string]string{"site": "s2"}
	if err := s.Update(context.Background(), OpsRecord{ID: "r1", Subject: "subject", Owner: "owner", Status: OpsStatusActive, Priority: OpsPriorityHigh, Labels: labels}, 0); err != nil {
		t.Fatal(err)
	}
	labels["site"] = "mutated"
	if got := s.items["r1"].Labels["site"]; got != "s2" {
		t.Fatalf("update aliased caller labels: %v", got)
	}
}
