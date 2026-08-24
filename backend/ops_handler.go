package main

import (
	"encoding/json"
	"io"
	"io/fs"
	"net/http"
	"strings"

	"example.com/hydro-gate-monitor-service/api"
	"example.com/hydro-gate-monitor-service/store"
)

// opsErrorBody is the structured body every ops error returns.
type opsErrorBody struct {
	Code      string `json:"code"`
	Message   string `json:"message"`
	Operation string `json:"operation,omitempty"`
	RequestID string `json:"request_id,omitempty"`
}

// opsSuccessResponse wraps every successful ops response payload.
type opsSuccessResponse struct {
	Data any `json:"data"`
}

// buildRouter mounts the ops endpoints under /api/ops and falls back to the
// existing gate/web router for everything else.
func buildRouter(s *store.Store, webFS fs.FS) http.Handler {
	ops := newOpsService(nil)
	mux := http.NewServeMux()
	mux.Handle("/api/ops", opsHandler(ops))
	mux.Handle("/api/ops/", opsHandler(ops))
	mux.Handle("/", api.NewRouter(s, webFS))
	return mux
}

func opsHandler(svc *OpsService) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Strip the "/api/ops" prefix, then route on the remainder.
		rest := strings.TrimPrefix(r.URL.Path, "/api/ops")
		rest = strings.Trim(rest, "/")
		parts := splitOpsPath(rest)

		switch {
		case len(parts) == 0:
			handleOpsCollection(w, r, svc)
		case len(parts) == 1:
			handleOpsRecord(w, r, svc, parts[0])
		case len(parts) == 2 && parts[1] == "audit":
			handleOpsAudit(w, r, svc, parts[0])
		default:
			opsWriteError(w, r, ErrOpsNotFound, "unknown ops subroute")
		}
	}
}

// splitOpsPath splits a path remainder into non-empty segments.
func splitOpsPath(rest string) []string {
	if rest == "" {
		return nil
	}
	out := []string{}
	for _, part := range strings.Split(rest, "/") {
		if part != "" {
			out = append(out, part)
		}
	}
	return out
}

func handleOpsCollection(w http.ResponseWriter, r *http.Request, svc *OpsService) {
	switch r.Method {
	case http.MethodPost:
		opsCreate(w, r, svc)
	case http.MethodGet:
		opsSearch(w, r, svc)
	default:
		opsWriteError(w, r, ErrOpsMethod, "method not allowed on collection")
	}
}

func handleOpsRecord(w http.ResponseWriter, r *http.Request, svc *OpsService, id string) {
	switch r.Method {
	case http.MethodGet:
		rec, err := svc.Get(r.Context(), id)
		if err != nil {
			opsWriteError(w, r, err, "get failed")
			return
		}
		opsWriteOK(w, http.StatusOK, rec)
	case http.MethodPatch:
		opsTransition(w, r, svc, id)
	case http.MethodDelete:
		if err := svc.Delete(r.Context(), id, opsActorFromRequest(r)); err != nil {
			opsWriteError(w, r, err, "delete failed")
			return
		}
		w.WriteHeader(http.StatusNoContent)
	default:
		opsWriteError(w, r, ErrOpsMethod, "method not allowed on record")
	}
}

func handleOpsAudit(w http.ResponseWriter, r *http.Request, svc *OpsService, id string) {
	if r.Method != http.MethodGet {
		opsWriteError(w, r, ErrOpsMethod, "method not allowed on audit")
		return
	}
	opsWriteOK(w, http.StatusOK, svc.Audit(id))
}

func opsCreate(w http.ResponseWriter, r *http.Request, svc *OpsService) {
	var record OpsRecord
	if err := decodeOpsBody(r, &record); err != nil {
		opsWriteError(w, r, err, "invalid request body")
		return
	}
	created, err := svc.Create(r.Context(), record)
	if err != nil {
		opsWriteError(w, r, err, "create failed")
		return
	}
	opsWriteOK(w, http.StatusCreated, created)
}

func opsSearch(w http.ResponseWriter, r *http.Request, svc *OpsService) {
	q := OpsQuery{
		Subject:  r.URL.Query().Get("subject"),
		Status:   OpsStatus(r.URL.Query().Get("status")),
		Priority: OpsPriority(r.URL.Query().Get("priority")),
		Owner:    r.URL.Query().Get("owner"),
		Page:     atoiOps(r.URL.Query().Get("page")),
		PageSize: atoiOps(r.URL.Query().Get("page_size")),
	}
	page, err := svc.Search(r.Context(), q)
	if err != nil {
		opsWriteError(w, r, err, "search failed")
		return
	}
	opsWriteOK(w, http.StatusOK, page)
}

func opsTransition(w http.ResponseWriter, r *http.Request, svc *OpsService, id string) {
	var body struct {
		Target   OpsStatus `json:"target"`
		Revision int       `json:"revision"`
	}
	if err := decodeOpsBody(r, &body); err != nil {
		opsWriteError(w, r, err, "invalid request body")
		return
	}
	if !opsStatusValid(body.Target) {
		opsWriteError(w, r, ErrOpsValidation, "target is not a valid status")
		return
	}
	rec, err := svc.Transition(r.Context(), id, body.Revision, body.Target, opsActorFromRequest(r))
	if err != nil {
		opsWriteError(w, r, err, "transition failed")
		return
	}
	opsWriteOK(w, http.StatusOK, rec)
}

func decodeOpsBody(r *http.Request, dst any) error {
	decoder := json.NewDecoder(io.LimitReader(r.Body, 1<<20))
	if err := decoder.Decode(dst); err != nil {
		return ErrOpsValidation
	}
	return nil
}

func opsWriteOK(w http.ResponseWriter, status int, value any) {
	opsJSON(w, status, opsSuccessResponse{Data: value})
}

// opsWriteError classifies err and always writes a structured JSON body so
// conflict (409) and not-found (404) are distinguishable, never a bare 200
// or empty body.
func opsWriteError(w http.ResponseWriter, r *http.Request, err error, message string) {
	body := opsErrorBody{
		Code:      opsCode(err),
		Message:   message,
		Operation: opsOperation(err),
		RequestID: opsRequestID(r),
	}
	opsJSON(w, opsHTTPStatus(err), body)
}

func opsOperation(err error) string {
	var typed *OpsError
	if errorsAsOps(err, &typed) {
		return typed.Operation
	}
	return ""
}

// errorsAsOps is a thin wrapper so this file does not need to import errors
// directly (keeps the dependency surface tiny).
func errorsAsOps(err error, target **OpsError) bool {
	if err == nil {
		return false
	}
	for cur := err; cur != nil; {
		if oe, ok := cur.(*OpsError); ok {
			*target = oe
			return true
		}
		type unwrapper interface{ Unwrap() error }
		u, ok := cur.(unwrapper)
		if !ok {
			break
		}
		cur = u.Unwrap()
	}
	return false
}

func atoiOps(value string) int {
	n := 0
	for _, ch := range value {
		if ch < '0' || ch > '9' {
			return 0
		}
		n = n*10 + int(ch-'0')
	}
	return n
}
