package main

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
	"time"

	internalgrpc "github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/server/grpc"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// test stores the HTTP testing client preconfigured
// var http = baloo.New(os.Getenv("APP_HOST"))

func TestMain(t *testing.T) {
	host := os.Getenv("APP_HOST")

	httpClient := http.Client{}

	res, err := httpClient.Get(host + "/events") // nolint: noctx
	if err != nil {
		t.Errorf("Failed to get /events: %s", err)
	}
	buf := strings.Builder{}
	io.Copy(&buf, res.Body)
	res.Body.Close()
	require.Equal(t, "[]", buf.String())

	// Создадим новое событие
	body := `{
		"id": "4927aa58-a175-429a-a125-c04765597152",
		"title": "Test Event 01",
		"description": "Test Event Description 01",
		"date": "2021-12-20 12:30:00",
		"duration": 60,
		"user_id": "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee",
		"notify_before_seconds": 60
	}`
	bufWrite := bytes.NewBuffer([]byte(body))

	res, err = httpClient.Post(host+"/events", "application/json", bufWrite) // nolint: noctx
	if err != nil {
		t.Errorf("Failed to create /events: %s", err)
	}
	res.Body.Close()

	res, err = httpClient.Get(host + "/events") // nolint: noctx
	if err != nil {
		t.Errorf("Failed to get /events: %s", err)
	}
	buf = strings.Builder{}
	io.Copy(&buf, res.Body)
	res.Body.Close()

	bodyExpected := `[{"id":"4927aa58-a175-429a-a125-c04765597152","title":"Test Event 01","date":"2021-12-20 12:30:00","duration":60,"description":"Test Event Description 01","user_id":"b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee","notify_before_seconds":60}]` // nolint:lll
	require.Equal(t, bodyExpected, buf.String())

	// Обновим событие
	body = `{
		"title": "Test Event 01 UPD",
		"description": "Test Event Description 01 UPD",
		"date": "2021-12-20 12:30:30",
		"duration": 70,
		"user_id": "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee",
		"notify_before_seconds": 70
	}`
	bufWrite = bytes.NewBuffer([]byte(body))

	req, err := http.NewRequest("PUT", host+"/events/4927aa58-a175-429a-a125-c04765597152", bufWrite) // nolint: noctx
	if err != nil {
		t.Errorf("Failed to update event: %s", err)
	}

	res, err = httpClient.Do(req)
	if err != nil {
		t.Errorf("Failed to update event: %s", err)
	}
	res.Body.Close()

	res, err = httpClient.Get(host + "/events") // nolint: noctx
	if err != nil {
		t.Errorf("Failed to get /events: %s", err)
	}
	buf = strings.Builder{}
	io.Copy(&buf, res.Body)
	res.Body.Close()

	bodyExpected = `[{"id":"4927aa58-a175-429a-a125-c04765597152","title":"Test Event 01 UPD","date":"2021-12-20 12:30:30","duration":70,"description":"Test Event Description 01 UPD","user_id":"b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee","notify_before_seconds":70}]` // nolint:lll
	require.Equal(t, bodyExpected, buf.String())

	// Удалим событие
	req, err = http.NewRequest("DELETE", host+"/events/4927aa58-a175-429a-a125-c04765597152", nil) // nolint: noctx
	if err != nil {
		t.Errorf("Failed to delete event: %s", err)
	}

	res, err = httpClient.Do(req)
	if err != nil {
		t.Errorf("Failed to delete event: %s", err)
	}
	res.Body.Close()

	res, err = httpClient.Get(host + "/events") // nolint: noctx
	if err != nil {
		t.Errorf("Failed to get /events: %s", err)
	}
	buf = strings.Builder{}
	io.Copy(&buf, res.Body)
	res.Body.Close()

	bodyExpected = `[]`
	require.Equal(t, bodyExpected, buf.String())

	time.Sleep(time.Second * 10)
	// Проверим, что было уведомление
	logFileName := "/var/logs/app.log"
	content, err := os.ReadFile(logFileName)
	if err != nil {
		t.Errorf("Failed to read sender logs")
		t.FailNow()
	}
	contentString := string(content)

	fmt.Println(contentString)

	require.Contains(t, contentString, "4927aa58-a175-429a-a125-c04765597152")
}

var host = os.Getenv("APP_HOST_GRPC")

func TestMainGrpc(t *testing.T) {
	conn, err := grpc.Dial(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		t.Errorf("Failed to dial GRPC service: %s", err)
	}

	ctx := context.Background()

	client := internalgrpc.NewEventServiceClient(conn)

	event := internalgrpc.Event{
		Id:                  "4927aa58-a175-429a-a125-c04765597152",
		Title:               "Event",
		Description:         "Event Descr",
		Date:                "2021-12-20 12:30:30",
		Duration:            45,
		UserId:              "b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee",
		NotifyBeforeSeconds: 100,
	}
	_, err = client.Create(ctx, &event)
	assert.Nil(t, err)

	req := internalgrpc.EventListRequest{
		Date: "2021-12-01",
	}
	res, err := client.EventListMonth(ctx, &req)
	assert.Nil(t, err)
	assert.Len(t, res.GetEvents(), 1)
}
