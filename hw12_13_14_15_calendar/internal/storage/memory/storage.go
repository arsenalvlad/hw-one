package memorystorage

import (
	"context"
	"sync"

	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/app"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/model"
)

type Storage struct {
	data map[int]model.Event
	mu   *sync.RWMutex
}

func New() app.Storage {
	return &Storage{
		data: nil,
		mu:   &sync.RWMutex{},
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

func (s *Storage) AddEvent(_ context.Context, data model.Event) (*model.Event, error) {
	s.mu.Lock()
	data.ID = len(s.data) + 1
	s.data[data.ID] = data
	s.mu.Unlock()

	return &data, nil
}

func (s *Storage) UpdateEvent(_ context.Context, data model.Event) (*model.Event, error) {
	s.mu.Lock()
	s.data[data.ID] = data
	s.mu.Unlock()

	return &data, nil
}

func (s *Storage) ListEvent(_ context.Context) ([]*model.Event, error) {
	if len(s.data) == 0 {
		return nil, nil
	}

	s.mu.Lock()
	items := make([]*model.Event, 0, len(s.data))

	for _, item := range s.data {
		res := item
		items = append(items, &res)
	}
	s.mu.Unlock()

	return items, nil
}
