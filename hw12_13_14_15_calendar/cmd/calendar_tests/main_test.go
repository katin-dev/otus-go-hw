package main

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// test stores the HTTP testing client preconfigured
// var http = baloo.New(os.Getenv("APP_HOST"))

func TestMain(t *testing.T) {
	/*
		1. Дождаться когда сервис поднимится
		2. CRUDL для HTTP
		3. CRUDL для gRPC
		4. Хз как проверить нотификацию
	*/

	host := os.Getenv("APP_HOST")

	httpClient := http.Client{}

	res, err := httpClient.Get(host + "/events")
	if err != nil {
		t.Errorf("Failed to get /events: %s", err)
	}
	buf := strings.Builder{}
	io.Copy(&buf, res.Body)

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

	_, err = httpClient.Post(host+"/events", "application/json", bufWrite)
	if err != nil {
		t.Errorf("Failed to create /events: %s", err)
	}

	res, err = httpClient.Get(host + "/events")
	if err != nil {
		t.Errorf("Failed to get /events: %s", err)
	}
	buf = strings.Builder{}
	io.Copy(&buf, res.Body)

	bodyExpected := `[{"id":"4927aa58-a175-429a-a125-c04765597152","title":"Test Event 01","date":"2021-12-20 12:30:00","duration":60,"description":"Test Event Description 01","user_id":"b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee","notify_before_seconds":60}]`
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

	req, err := http.NewRequest("PUT", host+"/events/4927aa58-a175-429a-a125-c04765597152", bufWrite)
	if err != nil {
		t.Errorf("Failed to update event: %s", err)
	}

	_, err = httpClient.Do(req)
	if err != nil {
		t.Errorf("Failed to update event: %s", err)
	}

	res, err = httpClient.Get(host + "/events")
	if err != nil {
		t.Errorf("Failed to get /events: %s", err)
	}
	buf = strings.Builder{}
	io.Copy(&buf, res.Body)

	bodyExpected = `[{"id":"4927aa58-a175-429a-a125-c04765597152","title":"Test Event 01 UPD","date":"2021-12-20 12:30:30","duration":70,"description":"Test Event Description 01 UPD","user_id":"b6a4fbfa-a9b2-429c-b0c5-20915c84e9ee","notify_before_seconds":70}]`
	require.Equal(t, bodyExpected, buf.String())

	// Удалим событие
	req, err = http.NewRequest("DELETE", host+"/events/4927aa58-a175-429a-a125-c04765597152", nil)
	if err != nil {
		t.Errorf("Failed to delete event: %s", err)
	}

	_, err = httpClient.Do(req)
	if err != nil {
		t.Errorf("Failed to delete event: %s", err)
	}

	res, err = httpClient.Get(host + "/events")
	if err != nil {
		t.Errorf("Failed to get /events: %s", err)
	}
	buf = strings.Builder{}
	io.Copy(&buf, res.Body)

	bodyExpected = `[]`
	require.Equal(t, bodyExpected, buf.String())
}
