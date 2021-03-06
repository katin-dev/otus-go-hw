package tools

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/app"
)

func CreateUUID(str string) uuid.UUID {
	id, _ := uuid.Parse(str)
	return id
}

func CreateDate(str string) time.Time {
	dt, _ := time.Parse(time.RFC3339, str)
	return dt
}

func ExtractEventID(events []app.Event) []string {
	res := make([]string, 0, len(events))

	for _, e := range events {
		res = append(res, e.ID.String())
	}

	return res
}

func JSONRemarshalString(body string) string {
	bytes := []byte(body)
	var ifce interface{}
	err := json.Unmarshal(bytes, &ifce)
	if err != nil {
		return ""
	}

	output, err := json.Marshal(ifce)
	if err != nil {
		return ""
	}

	return string(output)
}
