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

func (a *App) AddEvent(ctx context.Context, title string) error {
	event := model.Event{}
	event.Title = title

	_, err := a.Storage.AddEvent(ctx, event)
	if err != nil {
		return fmt.Errorf("could not add event: %w", err)
	}
	return nil
}

// TODO
