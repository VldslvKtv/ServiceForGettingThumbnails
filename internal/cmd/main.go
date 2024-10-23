package main

import (
	thumbapp "getthumbnails/internal/app/thumb_app"
	"getthumbnails/internal/config"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	conf := config.MustLoad()

	log := setupLogger(conf.Env)

	log.Info("start app",
		slog.String("env", conf.Env),
		slog.Any("cfg", conf),
		slog.Int("port", conf.GRPC.Port),
	)

	// TODO: ПОКА В ПРОЦЕССЕ
	application := thumbapp.New(log, conf.GRPC.Port, conf.StoragePath)

	go application.GRPCSrv.MustRun()

	// Gracefull shutdown

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT) // ждем сигнала от ОС - и  висим на строке <-stop
	// пока отдельная горутина запущена с сервером
	check := <-stop

	log.Info("server stopped", slog.String("signal", check.String()))

	application.GRPCSrv.Stop()
}

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func setupLogger(env string) *slog.Logger { // логирование
	var log *slog.Logger

	switch env {
	case envLocal:
		log = slog.New( // показываем куда выводить и какой уровень логирования
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envDev:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	case envProd:
		log = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log

}
