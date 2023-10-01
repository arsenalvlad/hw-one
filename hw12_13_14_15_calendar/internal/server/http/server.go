package internalhttp

import (
	"context"
	"fmt"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/app"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/logger"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type Server struct { // TODO
	Address string
	*logger.Logger
	app.Application
	Server *http.Server
}

func NewServer(logger *logger.Logger, app app.Application, address string) *Server {
	router := http.NewServeMux()
	router.HandleFunc("/", myHandler)

	server := &http.Server{
		Addr:    address,
		Handler: router,
	}

	return &Server{
		Logger:      logger,
		Application: app,
		Address:     address,
		Server:      server,
	}
}

func (s *Server) Start(ctx context.Context) error {
	s.Logger.Info("Starting http server ",
		zap.String("DSN", s.Address),
	)

	loggingMiddleware := LoggingMiddleware(s.Logger)
	loggedRouter := loggingMiddleware(s.Server.Handler)

	err := http.ListenAndServe(s.Address, loggedRouter) //nolint: gosec
	if err != nil {
		return fmt.Errorf("could not listen and serve http: %w", err)
	}

	return nil
}

func myHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello!")
	time.Sleep(1 * time.Second)
}

func (s *Server) Stop(ctx context.Context) error {
	return s.Server.Shutdown(ctx)
}
