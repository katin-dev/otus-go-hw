package memory

import (
	"sync"

	"github.com/google/uuid"
	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/app"
)

type Storage struct {
	mu     sync.RWMutex
	events map[uuid.UUID]app.Event
}

func New() *Storage {
	return &Storage{
		events: make(map[uuid.UUID]app.Event),
	}
}

func (s *Storage) Create(e app.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[e.Id] = e

	return nil
}

func (s *Storage) Update(e app.Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[e.Id] = e

	return nil
}

func (s *Storage) Delete(id uuid.UUID) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, id)

	return nil
}

func (s *Storage) FindAll() ([]app.Event, error) {
	events := make([]app.Event, 0, len(s.events))
	for _, v := range s.events {
		events = append(events, v)
	}

	return events, nil
}
