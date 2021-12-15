package sqlstorage

import (
	"context"

	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/xtgo/uuid"
)

type Storage struct{}

func New() *Storage {
	return &Storage{}
}

func (s *Storage) Connect(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	// TODO
	return nil
}

func (s *Storage) Create(e app.Event) {
}

func (s *Storage) Update(e app.Event) {
}

func (s *Storage) Delete(id uuid.UUID) {
}

func (s *Storage) FindAll() []app.Event {
	events := make([]app.Event, 0)
	return events
}
