package app

import (
	"fmt"
	serv "getthumbnails/internal/grpc"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

type AppGRPC struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

func New(
	log *slog.Logger,
	loadService serv.Load,
	port int,
) *AppGRPC {

	gRPCServer := grpc.NewServer()

	serv.Register(gRPCServer, loadService)

	return &AppGRPC{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}
}

func (a *AppGRPC) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *AppGRPC) Run() error {
	const op = "grpcapp.Run"

	log := a.log.With(
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port)) // слушатель запросов
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("grpc server is runnig", slog.String("addr", l.Addr().String()))

	if err := a.gRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil

}

func (a *AppGRPC) Stop() {
	const op = "grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.gRPCServer.GracefulStop() // новые запросы блокируется - а старые дорабатываются

}
