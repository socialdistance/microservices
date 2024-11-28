package app

import (
	"context"
	"log/slog"
	"time"

	grpcapp "lib_isod_v2/file_service/internal/app/grpc"
	httpapp "lib_isod_v2/file_service/internal/app/http"
	httprouters "lib_isod_v2/file_service/internal/http/file_service"
	"lib_isod_v2/file_service/internal/services/file"
	"lib_isod_v2/file_service/internal/services/reader"
	"lib_isod_v2/file_service/internal/services/watcher"
	"lib_isod_v2/file_service/internal/storage/postgresql"
)

type App struct {
	GRPCServer  *grpcapp.App
	HTTPServer  *httpapp.Server
	Watcher     *watcher.Watcher
	FileService *file.File
}

func New(log *slog.Logger, grpcPort int, storagePath string, httpHost, httpPort string, tokenTTL time.Duration, createPath string, recoveryPath string) *App {
	storage, err := postgresql.New(context.Background(), storagePath)
	if err != nil {
		panic(err)
	}

	watcher, err := watcher.NewWatcher(log, createPath, recoveryPath)
	if err != nil {
		panic(err)
	}
	reader := reader.NewReader(log)

	file := file.New(log, storage, watcher, reader)

	// authService := auth.New(log, storage, storage, storage, tokenTTL)

	grpcApp := grpcapp.New(log /* authService,*/, grpcPort)

	httpRouters := httprouters.NewRouter(log, storage)
	httpApp := httpapp.New(log, httpHost, httpPort, httpRouters)

	return &App{
		GRPCServer:  grpcApp,
		Watcher:     watcher,
		HTTPServer:  httpApp,
		FileService: file,
	}
}
