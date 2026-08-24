package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"example.com/hydro-gate-monitor-service/store"
)

func TestErrHTTPMissingGateNotFound(t *testing.T) {
	router := NewRouter(store.New(), nil)
	server := httptest.NewServer(router)
	defer server.Close()
	resp, err := http.Post(server.URL+"/api/gates/gate-999/ack", "application/json", strings.NewReader(`{"note":"reviewed"}`))
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNotFound {
		t.Fatalf("missing gate ack should be 404, got %d", resp.StatusCode)
	}
}
