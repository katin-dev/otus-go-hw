package app_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/storage/memory"
	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/tools"
	"github.com/stretchr/testify/require"
)

func TestAppEventCrud(t *testing.T) {
	logFile, err := os.CreateTemp("", "log")
	if err != nil {
		t.Errorf("failed to open test log file: %s", err)
	}

	logger, err := logger.New(logFile.Name(), "debug", "text_simple")
	if err != nil {
		t.Errorf("failed to open test log file: %s", err)
	}

	inMemoryStorage := memory.New()

	ctx := context.Background()

	testApp := app.New(logger, inMemoryStorage)

	event := app.Event{
		ID:           tools.CreateUUID("4927aa58-a175-429a-a125-c04765597152"),
		Title:        "Test Event",
		Description:  "Test Event Description",
		Dt:           tools.CreateDate("2021-12-20T00:00:00Z"),
		Duration:     time.Hour,
		UserID:       "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee",
		NotifyBefore: time.Minute * 15,
	}
	err = testApp.CreateEvent(ctx, event)
	require.Nil(t, err)

	// + week
	event.ID = tools.CreateUUID("11237ae6-a6f7-432d-90ba-351f17510a00")
	event.Dt = tools.CreateDate("2021-12-26T23:59:59Z")
	err = testApp.CreateEvent(ctx, event)
	require.Nil(t, err)

	// + month
	event.ID = tools.CreateUUID("45aad0db-284a-42a4-b3b5-b525937c688f")
	event.Dt = tools.CreateDate("2022-01-19T23:59:59Z")
	err = testApp.CreateEvent(ctx, event)
	require.Nil(t, err)

	// - day
	event.ID = tools.CreateUUID("5d1473a4-2e09-4424-ba2f-6ce771bc433c")
	event.Dt = tools.CreateDate("2021-12-19T23:59:59Z")
	err = testApp.CreateEvent(ctx, event)
	require.Nil(t, err)

	events, err := testApp.GetEventsByDay(ctx, tools.CreateDate("2021-12-20T07:15:45Z"))
	require.Nil(t, err)
	require.Len(t, events, 1)
	require.Equal(t, "4927aa58-a175-429a-a125-c04765597152", events[0].ID.String())

	events, err = testApp.GetEventsByWeek(ctx, tools.CreateDate("2021-12-20T07:15:45Z"))
	require.Nil(t, err)
	require.Len(t, events, 2)
	require.Equal(t, "4927aa58-a175-429a-a125-c04765597152", events[0].ID.String())
	require.Equal(t, "11237ae6-a6f7-432d-90ba-351f17510a00", events[1].ID.String())

	events, err = testApp.GetEventsByMonth(ctx, tools.CreateDate("2021-12-20T07:15:45Z"))
	require.Nil(t, err)
	require.Len(t, events, 3)
	require.Equal(t, "4927aa58-a175-429a-a125-c04765597152", events[0].ID.String())
	require.Equal(t, "11237ae6-a6f7-432d-90ba-351f17510a00", events[1].ID.String())
	require.Equal(t, "45aad0db-284a-42a4-b3b5-b525937c688f", events[2].ID.String())
}
