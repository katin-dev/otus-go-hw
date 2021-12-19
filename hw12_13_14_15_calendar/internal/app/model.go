package app

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrDateBusy      error = errors.New("date is already occupied")
	ErrRequiredField error = errors.New("required field")
)

type Event struct {
	Id           uuid.UUID
	Title        string
	Dt           time.Time
	Duration     time.Duration
	Description  string
	UserId       string
	NotifyBefore time.Duration
}

func NewEvent(title string, dt time.Time, duration time.Duration, userId string) *Event {
	id, _ := uuid.NewRandom()
	return &Event{
		Id:       id,
		Title:    title,
		Dt:       dt,
		Duration: duration,
		UserId:   userId,
	}
}
