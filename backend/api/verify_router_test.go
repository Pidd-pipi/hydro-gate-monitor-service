package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"testing/fstest"

	"example.com/hydro-gate-monitor-service/store"
)

func TestAssetRouterConcurrentStaticNoRace(t *testing.T) {
	files := fstest.MapFS{}
	for j := 0; j < 30; j++ {
		name := fmt.Sprintf("page-%02d.html", j)
		files[name] = &fstest.MapFile{Data: []byte(fmt.Sprintf("<html>page %d</html>", j))}
	}
	router := NewRouter(store.New(), files)
	server := httptest.NewServer(router)
	defer server.Close()
	const workers = 8
	start := make(chan struct{})
	var wg sync.WaitGroup
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			<-start
			for j := 0; j < 30; j++ {
				name := fmt.Sprintf("page-%02d.html", (j+i*7)%30)
				resp, err := http.Get(server.URL + "/" + name)
				if err == nil {
					_, _ = resp.Body.Read(nil)
					_ = resp.Body.Close()
				}
			}
		}(i)
	}
	close(start)
	wg.Wait()
}

func TestAssetStaticServesIndexAndAppJS(t *testing.T) {
	router := NewRouter(store.New(), fstest.MapFS{
		"page.html": {Data: []byte("<html>page</html>")},
		"script.js": {Data: []byte("console.log(1)")},
	})
	server := httptest.NewServer(router)
	defer server.Close()
	resp, err := http.Get(server.URL + "/page.html")
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("page status: %d", resp.StatusCode)
	}
	resp2, err2 := http.Get(server.URL + "/script.js")
	if err2 != nil {
		t.Fatal(err2)
	}
	defer resp2.Body.Close()
	if resp2.StatusCode != http.StatusOK {
		t.Fatalf("script status: %d", resp2.StatusCode)
	}
	if !strings.Contains(resp2.Header.Get("Content-Type"), "javascript") {
		t.Fatalf("script content type wrong: %q", resp2.Header.Get("Content-Type"))
	}
}
