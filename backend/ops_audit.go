package main

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

// opsAuditMaxEvents bounds the in-memory audit log so it cannot grow without
// limit over the lifetime of the process. When full, the oldest entries are
// evicted as new ones are appended.
const opsAuditMaxEvents = 10000

var opsAuditSequence uint64

func newOpsAuditID() string { return fmt.Sprintf("evt-%06d", atomic.AddUint64(&opsAuditSequence, 1)) }

type OpsAudit struct {
	mu     sync.RWMutex
	events []OpsEvent
}

func newOpsAudit() *OpsAudit { return &OpsAudit{events: make([]OpsEvent, 0, opsAuditMaxEvents)} }

// Add appends an audit event for the given record. The returned event is a
// stable copy of what was recorded.
func (a *OpsAudit) Add(recordID, typ, actor string) OpsEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	event := OpsEvent{ID: newOpsAuditID(), RecordID: recordID, Type: typ, Actor: actor, At: time.Now().UTC().Format(time.RFC3339Nano)}
	a.events = append(a.events, event)
	// Evict the oldest entries once the bound is exceeded so the slice cannot
	// grow unboundedly over the lifetime of the process.
	if len(a.events) > opsAuditMaxEvents {
		a.events = a.events[len(a.events)-opsAuditMaxEvents:]
		// Reclaim the backing array head so evicted entries do not pin memory.
		copy(a.events, a.events)
		a.events = a.events[:len(a.events)]
	}
	return event
}
func (a *OpsAudit) For(recordID string) []OpsEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := []OpsEvent{}
	for _, event := range a.events {
		if event.RecordID == recordID {
			out = append(out, event)
		}
	}
	return out
}
func (a *OpsAudit) Since(start time.Time) []OpsEvent {
	a.mu.RLock()
	defer a.mu.RUnlock()
	out := []OpsEvent{}
	for _, event := range a.events {
		parsed, err := time.Parse(time.RFC3339Nano, event.At)
		if err == nil && !parsed.Before(start) {
			out = append(out, event)
		}
	}
	return out
}
func (a *OpsAudit) Count() int { a.mu.RLock(); defer a.mu.RUnlock(); return len(a.events) }
func (a *OpsAudit) Latest() (OpsEvent, bool) {
	a.mu.RLock()
	defer a.mu.RUnlock()
	if len(a.events) == 0 {
		return OpsEvent{}, false
	}
	return a.events[len(a.events)-1], true
}
func (a *OpsAudit) Clear() { a.mu.Lock(); defer a.mu.Unlock(); a.events = a.events[:0] }
