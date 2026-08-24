package health

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestAssetHealthConcurrentNoRace(t *testing.T) {
	handler := http.HandlerFunc(Handler)
	server := httptest.NewServer(handler)
	defer server.Close()
	const workers = 8
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 50; j++ {
				resp, err := http.Get(server.URL + "/healthz")
				if err == nil {
					_ = resp.Body.Close()
				}
			}
		}()
	}
	close(start)
	wg.Wait()
}

func TestAssetHealthRejectsNonGET(t *testing.T) {
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/healthz", nil)
	Handler(rec, req)
	if rec.Code != http.StatusMethodNotAllowed {
		t.Fatalf("health should reject POST with 405, got %d", rec.Code)
	}
}
