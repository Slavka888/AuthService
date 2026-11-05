package app

import (
	grpcApp "AuthService/internal/app/grpc"
	authservice_service "AuthService/internal/services/auth-service"
	"AuthService/internal/storage/sqlite"
	"fmt"
	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcApp.App
	storage    *sqlite.Storage
}

func NewApp(
	log *slog.Logger,
	port int,
	storagePath string,
	tokenTTL time.Duration,
) (*App, error) {
	storage, err := sqlite.New(storagePath, log)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize storage: %w", err)
	}

	authService := authservice_service.New(log, storage, storage, storage, tokenTTL)

	grpcServer := grpcApp.NewApp(log, port, authService)

	return &App{
		GRPCServer: grpcServer,
		storage:    storage,
	}, nil
}

func (a *App) Stop() error {
	a.GRPCServer.Stop()
	return a.storage.Close()
}
