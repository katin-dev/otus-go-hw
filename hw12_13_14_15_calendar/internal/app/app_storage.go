package app

import (
	"github.com/google/uuid"
)

type Storage interface {
	Create(e Event) error
	Update(e Event) error
	Delete(id uuid.UUID) error
	FindAll() ([]Event, error)
}
