package grpc

import (
	"context"
	"fmt"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/logger"
	"github.com/arsenalvlad/hw12_13_14_15_calendar/internal/server/grpc/calendar"
	"go.uber.org/zap"
	"net"

	"google.golang.org/grpc"
)

func NewPort(dsn string, logger *logger.Logger, server calendar.ApplicationServer) (*Port, error) {
	return &Port{
		dsn:    dsn,
		logger: logger,
		server: server,
	}, nil
}

type Port struct {
	dsn    string
	logger *logger.Logger
	port   *grpc.Server
	server calendar.ApplicationServer
}

func loggedGRPCUnary(logger *logger.Logger) func(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler) (interface{}, error) {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		logger.Info("unary interceptor", zap.Any("Unary", info.FullMethod), zap.Any("Objects", req))
		return handler(ctx, req)
	}
}

func loggedGRPCStream(logger *logger.Logger) func(
	srv interface{},
	stream grpc.ServerStream,
	info *grpc.StreamServerInfo,
	handler grpc.StreamHandler) error {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		logger.Info("stream interceptor", zap.Any("Stream", info.FullMethod))
		return handler(srv, stream)
	}
}

func (p Port) Start() error {
	listener, err := net.Listen("tcp", p.dsn)
	if err != nil {
		return fmt.Errorf("could not start listener: %w", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(loggedGRPCUnary(p.logger)),
		grpc.StreamInterceptor(loggedGRPCStream(p.logger)))

	p.port = grpcServer

	p.logger.Info("Starting grpc server", zap.Any("DSN", p.dsn))

	calendar.RegisterApplicationServer(p.port, p.server)

	err = p.port.Serve(listener)
	if err != nil {
		return fmt.Errorf("could not serve listener: %w", err)
	}

	return nil
}

func (p Port) Stop() {
	//p.port.GracefulStop()
	p.logger.Info("grpc server stopped")
	p.port.Stop()
}
