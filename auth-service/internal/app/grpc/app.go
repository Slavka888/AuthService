package grpcApp

import (
	authservice "AuthService/internal/grpc/auth-service"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func NewApp(
	log *slog.Logger,
	port int,
	authService authservice.Auth,
) *App {

	gRPCServer := grpc.NewServer()

	authservice.Register(gRPCServer, authService)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *App) Run() error {
	log := a.log.With(
		slog.String("operation", "grpcapp.Run"),
		slog.Int("port", a.port),
	)

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		return fmt.Errorf("grpcapp.Run: %w", err)
	}

	log.Info("gRPC server is running", slog.Int("port", a.port))

	if err := a.gRPCServer.Serve(listener); err != nil {
		return fmt.Errorf("grpcapp.Run: %w", err)
	}

	return nil
}

func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *App) Stop() {
	log := a.log.With(
		slog.String("operation", "grpcapp.Stop"),
		slog.Int("port", a.port),
	)
	log.Info("Stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop()
}
