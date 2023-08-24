package sqlstorage

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/storage/sql/serializers"
	"github.com/volatiletech/null/v8"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/app"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/model"
	sqlModel "github.com/arsenalvlad/hw12_13_14_15_calendar/internal/storage/sql/model"
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
	var item *sqlModel.CalendarEvent

	item.Title = data.Title
	item.Duration = int64(data.Duration)
	item.EventTime = null.TimeFromPtr(&data.EventTime)
	item.Description = data.Description
	item.UserID = int64(data.UserID)

	err := item.Insert(ctx, s.db, boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("could not insert calendar event: %w", err)
	}

	return serializers.ToModelEvent(item), nil
}

func (s *Storage) UpdateEvent(ctx context.Context, data model.Event) (*model.Event, error) {
	item, err := sqlModel.FindCalendarEvent(ctx, s.db, int64(data.ID))
	if err != nil {
		return nil, fmt.Errorf("could not find calendar event: %w", err)
	}

	item.Title = data.Title
	item.Duration = int64(data.Duration)
	item.EventTime = null.TimeFromPtr(&data.EventTime)
	item.Description = data.Description
	item.UserID = int64(data.UserID)

	err = item.Upsert(ctx, s.db, false, nil, boil.Infer(), boil.Infer())
	if err != nil {
		return nil, fmt.Errorf("could not upsert calendar event: %w", err)
	}

	return serializers.ToModelEvent(item), nil
}

func (s *Storage) ListEvent(ctx context.Context) ([]*model.Event, error) {
	items, err := sqlModel.CalendarEvents().All(ctx, s.db)
	if err != nil {
		return nil, fmt.Errorf("could not get list calendar events: %w", err)
	}

	return serializers.ToModelEventSlice(items), nil
}
