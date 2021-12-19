package logger

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	now := time.Now()
	nowStr := now.Format(time.RFC3339)

	t.Run("test dev", func(t *testing.T) {
		file, err := os.CreateTemp("/tmp", "log")
		if err != nil {
			t.FailNow()
			return
		}

		defer os.Remove(file.Name())
		defer file.Close()

		l, _ := New(file.Name(), "dev")
		l.Debug("DEBUG", "param", 1)
		l.Info("INFO", "param", 2)
		// WARN & ERROR печатают TRACE потому их сложно проверять
		l.Flush()

		logContent, _ := os.ReadFile(file.Name())
		logExpected := nowStr + "\tDEBUG\tlogger/logger.go:52\tDEBUG\t{\"param\": 1}\n" +
			nowStr + "\tINFO\tlogger/logger.go:56\tINFO\t{\"param\": 2}\n"

		require.Equal(t, logExpected, string(logContent))
	})

	t.Run("test prod", func(t *testing.T) {
		file, err := os.CreateTemp("/tmp", "log")
		if err != nil {
			t.FailNow()
			return
		}

		defer os.Remove(file.Name())
		defer file.Close()

		l, _ := New(file.Name(), "prod")
		l.Debug("DEBUG", "param", 1)
		l.Info("INFO", "param", 2)
		l.Flush()

		// file.Close()
		logContent, _ := os.ReadFile(file.Name())
		fmt.Println(string(logContent))
		logExpected := "{\"level\":\"info\",\"timestamp\":\"" + nowStr + "\",\"caller\":\"logger/logger.go:56\",\"msg\":\"INFO\",\"param\":2}\n"

		require.Equal(t, logExpected, string(logContent))
	})
}
