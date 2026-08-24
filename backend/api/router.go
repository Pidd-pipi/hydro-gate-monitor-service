package api

import (
	"io/fs"
	"net/http"
	"path"
	"strings"
	"sync"

	"example.com/hydro-gate-monitor-service/health"
	"example.com/hydro-gate-monitor-service/store"
)

var (
	assetCache   = map[string][]byte{}
	assetCacheMu sync.RWMutex
)

func NewRouter(s *store.Store, webFS fs.FS) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", health.Handler)
	mux.HandleFunc("/api/gates", listGates(s))
	mux.HandleFunc("/api/gates/", changeGate(s))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		asset := r.URL.Path
		if asset == "/" {
			asset = "/index.html"
		}
		key := strings.TrimPrefix(asset, "/")
		assetCacheMu.RLock()
		data, ok := assetCache[key]
		assetCacheMu.RUnlock()
		if ok {
			writeAsset(w, asset, data)
			return
		}
		data, err := fs.ReadFile(webFS, path.Clean(key))
		if err != nil {
			http.NotFound(w, r)
			return
		}
		assetCacheMu.Lock()
		assetCache[key] = data
		assetCacheMu.Unlock()
		writeAsset(w, asset, data)
	})
	return mux
}

func writeAsset(w http.ResponseWriter, asset string, data []byte) {
	if strings.HasSuffix(asset, ".js") {
		w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	} else {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
	}
	_, _ = w.Write(data)
}

func pathID(path string) string {
	return strings.TrimSuffix(strings.TrimPrefix(path, "/api/gates/"), "/ack")
}
