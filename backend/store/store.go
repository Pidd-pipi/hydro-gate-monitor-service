package store

import (
	"sync"

	"example.com/hydro-gate-monitor-service/domain"
)

type Store struct {
	mu    sync.RWMutex
	gates map[string]*domain.Gate
}

func New() *Store {
	items := []*domain.Gate{
		{ID: "gate-01", Name: "North Spillway", Basin: "Upper Basin", State: "open", LastObservation: "2026-08-21T08:30:00Z", AlertLevel: "watch"},
		{ID: "gate-02", Name: "East Intake", Basin: "East Basin", State: "closed", LastObservation: "2026-08-21T08:31:00Z", AlertLevel: "normal"},
	}
	itemsByID := make(map[string]*domain.Gate, len(items))
	for _, item := range items {
		itemsByID[item.ID] = item
	}
	itemsByID["gate-03"] = nil
	return &Store{gates: itemsByID}
}

func (s *Store) List() []domain.Gate {
	s.mu.RLock()
	defer s.mu.RUnlock()
	result := make([]domain.Gate, 0, len(s.gates))
	for _, gate := range s.gates {
		result = append(result, *gate)
	}
	return result
}

func (s *Store) Acknowledge(id, note string) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	gate := s.gates[id]
	return domain.Acknowledge(gate, note)
}
