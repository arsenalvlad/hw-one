package memorystorage

import (
	"context"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/app"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/model"
	"sync"
)

type Storage struct {
	data map[int]model.Event
	mu   sync.RWMutex //nolint:unused
}

func New() app.Storage {
	return &Storage{
		data: nil,
		mu:   sync.RWMutex{},
	}
}

func (s *Storage) Connect() error {
	s.data = make(map[int]model.Event)

	return nil
}

func (s *Storage) Close() error {
	s.data = nil

	return nil
}

func (s *Storage) AddEvent(ctx context.Context, data model.Event) (*model.Event, error) {
	s.mu.Lock()
	data.ID = len(s.data) + 1
	s.data[data.ID] = data
	s.mu.Unlock()

	return &data, nil
}

func (s *Storage) UpdateEvent(ctx context.Context, data model.Event) (*model.Event, error) {
	s.mu.Lock()
	s.data[data.ID] = data
	s.mu.Unlock()

	return &data, nil
}

func (s *Storage) ListEvent(ctx context.Context) ([]*model.Event, error) {
	var items []*model.Event

	s.mu.Lock()
	for _, item := range s.data {
		items = append(items, &item)
	}
	s.mu.Unlock()

	return items, nil
}
