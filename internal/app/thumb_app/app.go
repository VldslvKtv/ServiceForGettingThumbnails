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
	//authService := auth.New(log, storage, storage, storage, tokenTTL)
	// TODO: инициализировать хранилище

	// TODO: init suth service (auth)

	grpcApp := app.New(log, thumbnailService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}
