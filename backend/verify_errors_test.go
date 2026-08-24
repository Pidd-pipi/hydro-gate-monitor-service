package main

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrWrapPreservesOpsErrorType(t *testing.T) {
	err := wrapOps("create", "store.put", ErrOpsNotFound)
	var typed *OpsError
	if !errors.As(err, &typed) {
		t.Fatalf("wrapOps lost *OpsError type: %v", err)
	}
	if typed.Code != "create" || typed.Operation != "store.put" {
		t.Fatalf("wrapOps wrong fields: %+v", typed)
	}
}

func TestErrCodeClassifiesWrappedErrors(t *testing.T) {
	chain := fmt.Errorf("outer: %w", ErrOpsConflict)
	if got := opsCode(chain); got != "conflict" {
		t.Fatalf("opsCode misclassified wrapped conflict: %q", got)
	}
	notFound := fmt.Errorf("outer: %w", ErrOpsNotFound)
	if got := opsCode(notFound); got != "not_found" {
		t.Fatalf("opsCode misclassified wrapped not-found: %q", got)
	}
}

func TestErrIsNotFoundMatchesWrappedError(t *testing.T) {
	chain := fmt.Errorf("outer: %w", ErrOpsNotFound)
	if !opsIsNotFound(chain) {
		t.Fatalf("opsIsNotFound failed on wrapped sentinel: %v", chain)
	}
}

func TestErrIsConflictMatchesWrappedError(t *testing.T) {
	chain := fmt.Errorf("outer: %w", ErrOpsConflict)
	if !opsIsConflict(chain) {
		t.Fatalf("opsIsConflict failed on wrapped sentinel: %v", chain)
	}
}

func TestErrIsInvalidMatchesWrappedError(t *testing.T) {
	chain := fmt.Errorf("outer: %w", ErrOpsInvalid)
	if !opsIsInvalid(chain) {
		t.Fatalf("opsIsInvalid failed on wrapped sentinel: %v", chain)
	}
}

func TestErrIsTransitionMatchesWrappedError(t *testing.T) {
	chain := fmt.Errorf("outer: %w", ErrOpsTransition)
	if !opsIsTransition(chain) {
		t.Fatalf("opsIsTransition failed on wrapped sentinel: %v", chain)
	}
}

func TestErrIsPolicyMatchesWrappedError(t *testing.T) {
	chain := fmt.Errorf("outer: %w", ErrOpsPolicy)
	if !opsIsPolicy(chain) {
		t.Fatalf("opsIsPolicy failed on wrapped sentinel: %v", chain)
	}
}

