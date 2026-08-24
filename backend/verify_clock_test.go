package main

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestCtxOpsContextKeepsParentDeadline(t *testing.T) {
	parent, cancel := context.WithTimeout(context.Background(), time.Minute)
	defer cancel()
	ctx, ccancel := opsContext(parent, 3*time.Second)
	defer ccancel()
	select {
	case <-ctx.Done():
		t.Fatal("derived context should not be done before parent")
	default:
	}
	parentCancel, pcancel := context.WithCancel(context.Background())
	pcancel()
	cctx, cccancel := opsContext(parentCancel, 3*time.Second)
	defer cccancel()
	select {
	case <-cctx.Done():
	case <-time.After(500 * time.Millisecond):
		t.Fatal("parent cancellation was not propagated into opsContext")
	}
}

func TestCtxOpsDelayHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(50 * time.Millisecond)
		cancel()
	}()
	start := time.Now()
	err := opsDelay(ctx, 2*time.Second)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("delay ignored cancellation: %v", err)
	}
	if time.Since(start) > time.Second {
		t.Fatalf("delay was not interrupted: %v", time.Since(start))
	}
}

func TestCtxTransitionHonorsParentCancel(t *testing.T) {
	svc := newOpsService([]OpsRecord{{ID: "r1", Subject: "subject", Owner: "owner", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}}})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := svc.Transition(ctx, "r1", 0, OpsStatusActive, "alice")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("transition did not honor parent cancel: %v", err)
	}
}

func TestCtxServiceGetHonorsParentCancel(t *testing.T) {
	svc := newOpsService([]OpsRecord{{ID: "r1", Subject: "subject", Owner: "owner", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}}})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := svc.Get(ctx, "r1")
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("service get did not honor parent cancel: %v", err)
	}
}

func TestCtxCreateHonorsParentCancel(t *testing.T) {
	svc := newOpsService(nil)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := svc.Create(ctx, OpsRecord{ID: "r9", Subject: "subject", Owner: "owner", Status: OpsStatusQueued, Priority: OpsPriorityNormal, Labels: map[string]string{"site": "s1"}})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("create did not honor parent cancel: %v", err)
	}
}
