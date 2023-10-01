package grpc

import (
	"context"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/app"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/server/grpc/calendar"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ calendar.ApplicationServer = (*Server)(nil)

func NewServer(app *app.App) (*Server, error) {
	return &Server{app: app}, nil
}

type Server struct {
	app *app.App
}

func (s Server) AddEvent(ctx context.Context, event *calendar.Event) (*calendar.Event, error) {
	//TODO implement me
	panic("implement me")
}

func (s Server) UpdateEvent(ctx context.Context, event *calendar.Event) (*calendar.Event, error) {
	//TODO implement me
	return &calendar.Event{
		Id:          0,
		Title:       "ssadasdasd",
		Description: "",
		UserId:      0,
		EventTimeAt: nil,
		Duration:    nil,
	}, nil
}

func (s Server) ListEvents(ctx context.Context, empty *emptypb.Empty) (*calendar.EventSlice, error) {
	//TODO implement me
	return nil, nil
}
