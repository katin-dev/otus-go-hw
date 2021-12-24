package sql

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4"
	"github.com/katin.dev/otus-go-hw/hw12_13_14_15_calendar/internal/app"
)

type Storage struct {
	ctx  context.Context
	conn *pgx.Conn
	dsn  string
}

func New(ctx context.Context, dsn string) *Storage {
	return &Storage{
		dsn: dsn,
		ctx: ctx,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, s.dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	s.conn = conn

	return nil
}

func (s *Storage) Close(ctx context.Context) error {
	return s.conn.Close(ctx)
}

func (s *Storage) Create(e app.Event) error {
	var sql string
	/* sql := "SELECT * FROM events WHERE date = $1 AND user_id = $1"
	row := s.conn.QueryRow(s.ctx, sql, e.Dt.Unix(), e.UserId)
	if err := row.Scan(); err != nil {
		if err != pgx.ErrNoRows {
			// return fmt.Errorf("event with such date and user already exists")
			return fmt.Errorf("Failed to SELECT: %w", err)
		}
	} */

	sql = "INSERT INTO events(id, title, date, duration, description, user_id, notify_before) VALUES ($1, $2, $3, $4, $5, $6, $7)"
	_, err := s.conn.Exec(s.ctx, sql, e.Id.String(), e.Title, e.Dt.Format(time.RFC3339), e.Duration.Seconds(), e.Description, e.UserId, e.NotifyBefore.Seconds())

	return err
}

func (s *Storage) Update(e app.Event) error {
	sql := "UPDATE events SET title=$1, date = $2, duration=$3, description=$4, user_id=$5, notify_before=$6 WHERE id = $7"
	_, err := s.conn.Exec(s.ctx, sql, e.Title, e.Dt.Format(time.RFC3339), e.Duration.Seconds(), e.Description, e.UserId, e.NotifyBefore.Seconds(), e.Id.String())

	return err
}

func (s *Storage) Delete(id uuid.UUID) error {
	sql := "DELETE FROM events WHERE id = $1"
	_, err := s.conn.Exec(s.ctx, sql, id)

	return err
}

func (s *Storage) FindAll() ([]app.Event, error) {
	events := make([]app.Event, 0)

	sql := "SELECT id, title, date, duration, description, user_id, notify_before FROM events ORDER BY date"
	rows, err := s.conn.Query(s.ctx, sql)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e app.Event
		var durationSeconds, notifyBeforeSeconds int
		if err := rows.Scan(&e.Id, &e.Title, &e.Dt, &durationSeconds, &e.Description, &e.UserId, &notifyBeforeSeconds); err != nil {
			return nil, fmt.Errorf("failed to scan SQL result into struct: %w", err)
		}

		e.Duration = time.Duration(durationSeconds * 1000000000)
		e.NotifyBefore = time.Duration(notifyBeforeSeconds * 1000000000)

		events = append(events, e)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}
