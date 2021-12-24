package sql

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/app"
	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

const APP_CONF_FILE = "configs/config.yaml"

func TestStorage(t *testing.T) {
	// Этот тест чисто для себя, в ci гонять его не надо:
	if _, err := os.Stat(APP_CONF_FILE); errors.Is(err, os.ErrNotExist) {
		t.Skip(APP_CONF_FILE + " file does not exists")
	}

	cfgContent, _ := os.ReadFile(APP_CONF_FILE)
	var cfg struct {
		Storage struct {
			Dsn string
		}
	}

	yaml.Unmarshal(cfgContent, cfg)

	ctx := context.Background()

	s := New(ctx, cfg.Storage.Dsn)
	if err := s.Connect(ctx); err != nil {
		t.Fatal("Failed to connect to DB server", err)
	}

	t.Run("test SQLStorage CRUDL", func(t *testing.T) {
		tx, _ := s.conn.BeginTx(ctx, pgx.TxOptions{
			IsoLevel:       pgx.Serializable,
			AccessMode:     pgx.ReadWrite,
			DeferrableMode: pgx.NotDeferrable,
		})

		userId := "3b6394b3-acc6-4fd5-b8ce-3cbdf30745ef"
		dt, _ := time.Parse("2006-01-02 15:04:05", "2021-01-01 12:00:00")

		event := app.NewEvent("Test Event", dt, time.Minute*30, userId)
		event.Description = "OTUS GoLang Lesson"
		event.NotifyBefore = time.Minute * 15

		err := s.Create(*event)
		require.Nil(t, err)

		saved, err := s.FindAll()
		require.Nil(t, err)
		require.Len(t, saved, 1)
		require.Equal(t, *event, saved[0])

		// Обновим параметры события:
		event.Title = "Test Event Upd"
		event.Description = "OTUS GoLang Lesson Upd"
		event.NotifyBefore = time.Minute * 15

		// Убедимся, что объект не изменился в репозитории только потому, что там хранятся ссылки,а не копии
		saved, _ = s.FindAll()
		require.Len(t, saved, 1)
		require.NotEqual(t, *event, saved[0])

		// Обновляем объект в репозитории
		err = s.Update(*event)
		if err != nil {
			t.Fatalf("Update failed: %s", err)
		}

		// Теперь он должен быть изменён
		saved, err = s.FindAll()
		if err != nil {
			t.Fatalf("failed to findAll after update: %s", err)
		}
		require.Len(t, saved, 1)
		require.Equal(t, *event, saved[0])

		// Удаляем объект
		s.Delete(event.Id)

		saved, _ = s.FindAll()
		require.Len(t, saved, 0)

		tx.Rollback(ctx)
	})
}
