package app

import (
	"context"
	"fmt"

	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/model"
	"go.uber.org/zap"
)

type App struct {
	Logger
	Storage
}

type Application interface { // TODO
	AddEvent(ctx context.Context, title string) (*model.Event, error)
}

type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	DPanic(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Sync() error
}

type Storage interface { // TODO
	Connect() error
	Close() error
	AddEvent(ctx context.Context, data model.Event) (*model.Event, error)
	UpdateEvent(ctx context.Context, data model.Event) (*model.Event, error)
	ListEvent(ctx context.Context) ([]*model.Event, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) AddEvent(ctx context.Context, title string) (*model.Event, error) {
	event := model.Event{}
	event.Title = title

	res, err := a.Storage.AddEvent(ctx, event)
	if err != nil {
		return nil, fmt.Errorf("could not add event: %w", err)
	}

	return res, nil
}
