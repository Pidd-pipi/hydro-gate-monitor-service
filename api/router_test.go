package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"

	"example.com/hydro-gate-monitor-service/store"
)

func testRouter() http.Handler {
	return NewRouter(store.New(), fstest.MapFS{"index.html": {Data: []byte("<html>hydro</html>")}, "app.js": {Data: []byte("fetch('/api/gates')")}})
}

func TestHealthAndList(t *testing.T) {
	server := httptest.NewServer(testRouter())
	defer server.Close()
	for _, path := range []string{"/healthz", "/api/gates", "/"} {
		resp, err := http.Get(server.URL + path)
		if err != nil {
			t.Fatal(err)
		}
		if resp.StatusCode != http.StatusOK {
			t.Fatalf("%s: got %d", path, resp.StatusCode)
		}
		_ = resp.Body.Close()
	}
}

func TestAcknowledgeValidationAndSuccess(t *testing.T) {
	server := httptest.NewServer(testRouter())
	defer server.Close()
	bad, _ := http.Post(server.URL+"/api/gates/gate-01/ack", "application/json", strings.NewReader(`{"note":""}`))
	if bad.StatusCode != http.StatusBadRequest {
		t.Fatalf("bad status: %d", bad.StatusCode)
	}
	_ = bad.Body.Close()
	good, _ := http.Post(server.URL+"/api/gates/gate-01/ack", "application/json", strings.NewReader(`{"note":"Operator reviewed spillway reading"}`))
	if good.StatusCode != http.StatusOK {
		t.Fatalf("good status: %d", good.StatusCode)
	}
	_ = good.Body.Close()
}
