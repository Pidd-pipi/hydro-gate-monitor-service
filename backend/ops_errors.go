package main

import (
	"errors"
	"net/http"
)

var (
	ErrOpsNotFound   = errors.New("operations record not found")
	ErrOpsConflict   = errors.New("operations revision conflict")
	ErrOpsInvalid    = errors.New("operations request is invalid")
	ErrOpsTransition = errors.New("operations status transition is not allowed")
	ErrOpsPolicy     = errors.New("operations policy rejected the request")
	ErrOpsValidation = errors.New("operations request failed validation")
	ErrOpsMethod     = errors.New("operations method not allowed")
)

// knownOpsCodes is the closed set of category codes that may appear as
// OpsError.Code. An operation name (e.g. "get", "store.put") must never be
// stored as the code — it would otherwise mask the real category. opsCode
// only trusts OpsError.Code when it belongs to this set; otherwise it
// unwraps the error and classifies the underlying cause.
var knownOpsCodes = map[string]bool{
	"not_found":        true,
	"conflict":         true,
	"invalid":          true,
	"transition":       true,
	"policy":           true,
	"validation":       true,
	"method_not_allowed": true,
	"internal":         true,
}

type OpsError struct {
	Code      string
	Operation string
	Cause     error
}

func (e *OpsError) Error() string {
	if e.Cause == nil {
		return e.Code + ": " + e.Operation
	}
	return e.Code + ": " + e.Operation + ": " + e.Cause.Error()
}
func (e *OpsError) Unwrap() error { return e.Cause }

// wrapOps preserves the error chain: it returns a *OpsError that Unwraps to
// cause, so errors.Is/As keep walking through wrapped sentinels. (Using
// fmt.Errorf with %v would flatten the cause to a string and sever the
// chain, collapsing every wrapped sentinel to "internal".)
func wrapOps(code, operation string, cause error) error {
	return &OpsError{Code: code, Operation: operation, Cause: cause}
}

// opsCode maps an error onto a stable category code. It uses errors.Is so
// sentinels wrapped via fmt.Errorf("%w", ...) or wrapOps still classify
// correctly.
func opsCode(err error) string {
	if err == nil {
		return ""
	}
	var typed *OpsError
	if errors.As(err, &typed) {
		if knownOpsCodes[typed.Code] {
			return typed.Code
		}
		// An opaque code (likely an operation name passed where a category
		// was expected) must not mask the real category — classify the
		// underlying cause instead.
		if typed.Cause != nil {
			return opsCode(typed.Cause)
		}
	}
	switch {
	case errors.Is(err, ErrOpsNotFound):
		return "not_found"
	case errors.Is(err, ErrOpsConflict):
		return "conflict"
	case errors.Is(err, ErrOpsInvalid), errors.Is(err, ErrOpsValidation):
		return "invalid"
	case errors.Is(err, ErrOpsTransition):
		return "transition"
	case errors.Is(err, ErrOpsPolicy):
		return "policy"
	case errors.Is(err, ErrOpsMethod):
		return "method_not_allowed"
	default:
		return "internal"
	}
}

// opsHTTPStatus maps a category code to an HTTP status.
func opsHTTPStatus(err error) int {
	switch opsCode(err) {
	case "not_found":
		return http.StatusNotFound
	case "conflict", "transition":
		return http.StatusConflict
	case "invalid":
		return http.StatusBadRequest
	case "policy":
		return http.StatusUnprocessableEntity
	case "method_not_allowed":
		return http.StatusMethodNotAllowed
	default:
		return http.StatusInternalServerError
	}
}

func opsIsNotFound(err error) bool   { return errors.Is(err, ErrOpsNotFound) }
func opsIsConflict(err error) bool   { return errors.Is(err, ErrOpsConflict) }
func opsIsInvalid(err error) bool    { return errors.Is(err, ErrOpsInvalid) || errors.Is(err, ErrOpsValidation) }
func opsIsTransition(err error) bool { return errors.Is(err, ErrOpsTransition) }
func opsIsPolicy(err error) bool     { return errors.Is(err, ErrOpsPolicy) }
