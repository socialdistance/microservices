package app

import (
	grpcapp "lib_isod_v2/auth_service/internal/app/grpc"
	httpapp "lib_isod_v2/auth_service/internal/app/http"
	httprouter "lib_isod_v2/auth_service/internal/http"

	"lib_isod_v2/auth_service/internal/services/auth"
	"lib_isod_v2/auth_service/internal/storage/sqlite"

	"log/slog"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
	HTTPServer *httpapp.Server
}

func New(log *slog.Logger, httpHost, httpPort string, grpcPort int, storagePath string, tokenTTL time.Duration) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log, authService, grpcPort)

	httpRouters := httprouter.NewRouters(log, authService)
	httpApp := httpapp.New(log, httpHost, httpPort, httpRouters)

	return &App{
		GRPCServer: grpcApp,
		HTTPServer: httpApp,
	}
}
