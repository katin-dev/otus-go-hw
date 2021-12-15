package memorystorage

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

func (s *Storage) Create(e app.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[e.Id] = e
}

func (s *Storage) Update(e app.Event) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.events[e.Id] = e
}

func (s *Storage) Delete(id uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.events, id)
}

func (s *Storage) FindAll() []app.Event {
	events := make([]app.Event, 0, len(s.events))
	for _, v := range s.events {
		events = append(events, v)
	}

	return events
}
