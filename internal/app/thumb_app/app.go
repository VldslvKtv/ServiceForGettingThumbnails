package thumbapp

import (
	"getthumbnails/internal/app"
	"getthumbnails/internal/loader"
	"getthumbnails/storage"
	"log/slog"
)

type App struct {
	GRPCSrv *app.AppGRPC
}

func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
) *App {
	storage, err := storage.New(storagePath)
	if err != nil {
		panic(err)
	}

	thumbnailService := loader.New(log, storage)

	grpcApp := app.New(log, thumbnailService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
