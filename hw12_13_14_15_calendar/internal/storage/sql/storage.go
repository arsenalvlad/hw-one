package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/app"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/model"
)

type Storage struct { // TODO
	db  *sql.DB
	dsn string
}

func New(dsn string) app.Storage {
	return &Storage{
		db:  nil,
		dsn: dsn,
	}
}

func (s *Storage) Connect() error {
	conn, err := sql.Open("postgres", s.dsn)
	if err != nil {
		return fmt.Errorf("could not open postgresql connection: %w", err)
	}

	s.db = conn

	err = s.db.Ping()
	if err != nil {
		return fmt.Errorf("could not ping postgresql database: %w", err)
	}

	return nil
}

func (s *Storage) Close() error {
	if s.db == nil {
		return nil
	}

	err := s.db.Close()
	if err != nil {
		return fmt.Errorf("could not close connect postgresql database: %w", err)
	}

	return nil
}

func (s *Storage) AddEvent(ctx context.Context, data model.Event) (*model.Event, error) {
	var result model.Event

	query := `INSERT INTO otus.public.calendar_event(
		"user_id", "title", "event_time","duration", "description"
	) values($1, $2, $3, $4, $5)
	RETURNING *;`

	err := s.db.QueryRowContext(ctx, query,
		data.UserID, data.Title, data.EventTime, data.Duration, data.Description,
	).Scan(&result.ID, &result.UserID, &result.Title, &result.EventTime, &result.Duration, &result.Description)
	if err != nil {
		return nil, fmt.Errorf("could not add new calendar event: %w", err)
	}

	return &result, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, data model.Event) (*model.Event, error) {
	query := `UPDATE otus.public.calendar_event
	SET "user_id" = $2, "title" = $3, "event_time" = $4,"duration" = $5, "description" = $6
	WHERE id=$1;`

	_, err := s.db.ExecContext(ctx, query,
		data.ID, data.UserID, data.Title, data.EventTime, data.Duration, data.Description)
	if err != nil {
		return nil, fmt.Errorf("could not update calendar event: %w", err)
	}

	return &data, nil
}

func (s *Storage) ListEvent(ctx context.Context) ([]*model.Event, error) {
	var results []*model.Event
	var item model.Event

	query := `SELECT "id", "user_id", "title", "event_time","duration", "description"
	FROM otus.public.calendar_event`

	rows, err := s.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("could not get list calendar events: %w", err)
	}

	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&item.ID, &item.UserID, &item.Title, &item.EventTime, &item.Duration, &item.Description)
		if err != nil {
			return nil, fmt.Errorf("could not scan list calendar events: %w", err)
		}

		results = append(results, &item)
	}

	return results, nil
}
