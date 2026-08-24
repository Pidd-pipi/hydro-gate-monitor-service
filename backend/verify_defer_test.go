package main

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"sync"
	"syscall"
	"testing"
	"time"
)

func TestLeakAuditAddReturnsRealEvent(t *testing.T) {
	a := newOpsAudit()
	ev := a.Add("r1", "created", "alice")
	if ev.ID == "" || ev.RecordID != "r1" || ev.Type != "created" || ev.Actor != "alice" {
		t.Fatalf("Add returned zero/bad event: %+v", ev)
	}
}

func TestLeakAuditEventsBounded(t *testing.T) {
	a := newOpsAudit()
	for i := 0; i < 15000; i++ {
		a.Add("r", "t", "a")
	}
	if got := a.Count(); got > 10000 {
		t.Fatalf("audit events unbounded: count=%d", got)
	}
}

func TestLeakRequestIDsUniqueUnderParallelLoad(t *testing.T) {
	handler := requestIDMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	const workers = 8
	const rounds = 50
	start := make(chan struct{})
	var wg sync.WaitGroup
	ids := make(chan string, workers*rounds)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < rounds; j++ {
				rec := httptest.NewRecorder()
				handler.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/x", nil))
				ids <- rec.Header().Get("X-Request-ID")
			}
		}()
	}
	close(start)
	wg.Wait()
	close(ids)
	seen := map[string]bool{}
	for id := range ids {
		if seen[id] {
			t.Fatalf("duplicate request id %q", id)
		}
		seen[id] = true
	}
}

func runShutdownCycle(t *testing.T) {
	t.Helper()
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	addr := ln.Addr().String()
	if err := ln.Close(); err != nil {
		t.Fatal(err)
	}
	srv := newEnterpriseServer(addr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	done := make(chan struct{})
	go func() {
		_ = serveHTTP(srv)
		close(done)
	}()
	deadline := time.Now().Add(2 * time.Second)
	for {
		conn, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err == nil {
			_ = conn.Close()
			break
		}
		if time.Now().After(deadline) {
			t.Fatal("server never started listening")
		}
		time.Sleep(10 * time.Millisecond)
	}
	if err := syscall.Kill(os.Getpid(), syscall.SIGTERM); err != nil {
		t.Fatal(err)
	}
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("serveHTTP did not return after signal")
	}
}

func TestLeakServeHTTPNoGoroutineLeakAfterShutdown(t *testing.T) {
	runShutdownCycle(t)
	time.Sleep(200 * time.Millisecond)
	before := runtime.NumGoroutine()
	for i := 0; i < 3; i++ {
		runShutdownCycle(t)
	}
	time.Sleep(300 * time.Millisecond)
	after := runtime.NumGoroutine()
	fmt.Printf("goroutines before=%d after=%d\n", before, after)
	if after > before+1 {
		t.Fatalf("goroutine count grew after shutdowns: before=%d after=%d", before, after)
	}
}
