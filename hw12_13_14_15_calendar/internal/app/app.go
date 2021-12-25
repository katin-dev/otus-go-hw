package app

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
)

type App struct {
	logg Logger
	repo Storage
}

func New(logger Logger, storage Storage) *App {
	return &App{
		logg: logger,
		repo: storage,
	}
}

func (a *App) CreateEvent(ctx context.Context, evt Event) error {
	var (
		err  error
		prev *Event
	)
	a.logg.Debug("App.CreateEvent - create new event")

	// Проверим, что он уже существует
	if prev, err = a.repo.FindOne(evt.ID); err != nil {
		a.logg.Error("App.CreateEvent ERROR: %s", err)
		return err
	}

	if prev != nil {
		a.logg.Warn("App.CreateEvent.AlreadyExists: %s", evt.ID)
		return fmt.Errorf("validation error: event with such id already exists: %s", evt.ID)
	}

	// Если ещё нет с таким ID - создаём
	if err = a.repo.Create(evt); err != nil {
		a.logg.Error("App.CreateEvent ERROR: %s", err)
		return err
	}

	return nil
}

func (a *App) UpdateEvent(ctx context.Context, evt Event) error {
	var (
		err  error
		prev *Event
	)
	a.logg.Debug("App.UpdateEvent.Begin %s", evt.ID)

	// Проверим наличие
	if prev, err = a.repo.FindOne(evt.ID); err != nil {
		a.logg.Error("App.UpdateEvent ERROR: %s", err)
		return err
	}

	if prev == nil {
		a.logg.Warn("App.UpdateEvent.NotFound: %s", evt.ID)
		return fmt.Errorf("validation error: event %s not found", evt.ID)
	}

	// Если ещё нет с таким ID - создаём
	if err = a.repo.Update(evt); err != nil {
		a.logg.Error("App.UpdateEvent ERROR: %s", err)
		return err
	}

	return nil
}

func (a *App) DeleteEvent(ctx context.Context, id uuid.UUID) error {
	var (
		err  error
		prev *Event
	)
	a.logg.Debug("App.DeleteEvent.Begin %s", id)

	// Проверим наличие
	if prev, err = a.repo.FindOne(id); err != nil {
		a.logg.Error("App.DeleteEvent ERROR: %s", err)
		return err
	}

	if prev == nil {
		a.logg.Warn("App.DeleteEvent.NotFound: %s", id)
		return fmt.Errorf("validation error: event %s not found", id)
	}

	// Если ещё нет с таким ID - создаём
	if err = a.repo.Delete(prev.ID); err != nil {
		a.logg.Error("App.DeleteEvent ERROR: %s", err)
		return err
	}

	return nil
}

func (a *App) GetEvents(ctx context.Context) ([]Event, error) {
	return a.repo.FindAll()
}

func (a *App) GetEventsByDay(ctx context.Context, day time.Time) ([]Event, error) {
	return a.GetEventsByInterval(ctx, day, time.Hour*24)
}

func (a *App) GetEventsByWeek(ctx context.Context, day time.Time) ([]Event, error) {
	return a.GetEventsByInterval(ctx, day, time.Hour*24*7)
}

func (a *App) GetEventsByMonth(ctx context.Context, day time.Time) ([]Event, error) {
	return a.GetEventsByInterval(ctx, day, time.Hour*24*7*30)
}

func (a *App) GetEventsByInterval(ctx context.Context, day time.Time, interval time.Duration) ([]Event, error) {
	events := make([]Event, 0)

	// Приведём время ко дню
	day, _ = time.Parse("2006-01-02", day.Format("2006-01-02"))

	a.logg.Debug("Get Event List from %s, interval: %s", day, interval)

	items, err := a.repo.FindAll()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		// Я точно знаю, что day начинается с 00:00:00, если я прибавлю 7 дней, то будет
		diff := item.Dt.Sub(day)
		if diff >= 0 && diff < interval {
			fmt.Printf("%s + %s >= %s\n", day, interval, item.Dt)
			events = append(events, item)
		} else {
			fmt.Printf("%s + %s < %s\n", day, interval, item.Dt)
		}
	}

	return events, nil
}
