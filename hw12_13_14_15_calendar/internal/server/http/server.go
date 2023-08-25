package internalhttp

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"go.uber.org/zap"
)

type Server struct { // TODO
	Address string
	Logger
	Application
}

type Logger interface { // TODO
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	DPanic(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)
	Sync() error
}

type Application interface { // TODO
	AddEvent(ctx context.Context, title string) error
}

func NewServer(logger Logger, app Application, address string) *Server {
	return &Server{
		Logger:      logger,
		Application: app,
		Address:     address,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.Logger.Info("Start server ",
		zap.String("address", s.Address),
	)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		s.Logger.Info("hadle /")
		_, err := io.WriteString(w, "Hi Otus!\n")
		if err != nil {
			s.Logger.Error("could not handle function to /: " + err.Error())
		}
	})

	err := http.ListenAndServe(s.Address, nil) //nolint: gosec
	if err != nil {
		return fmt.Errorf("could not listen and serve http: %w", err)
	}

	go func() {
		<-ctx.Done()
		err := s.Stop(ctx)
		if err != nil {
			s.Logger.Error("could not stop http server: " + err.Error())
		}
	}()

	return nil
}

func (s *Server) Stop(_ context.Context) error {
	os.Exit(1)

	return nil
}
