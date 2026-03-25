package grpcapp

import (
	"ReservationsService/internal/core"
	resrvgrpc "ReservationsService/internal/transport/grpc/reserv"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log/slog"
	"net"
)

type App struct {
	log        *slog.Logger
	gRPCServer *grpc.Server
	port       int
}

type ReservServise interface {
	AddReservation(ctx context.Context, request *core.ReservationRequest) (*core.ReservationResponse, error)
	CloseReservation(ctx context.Context, request *core.ReturnRequest) (*core.ReturnResponse, error)
	CheckReservation(ctx context.Context, userId int) (*core.CheckReservResponse, error)
}

type BookClientService interface {
	CheckBookAvailability(ctx context.Context, req *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error)
}

func New(log *slog.Logger, port int, reservService ReservServise, bookClient BookClientService) *App {
	gRPCServer := grpc.NewServer()

	resrvgrpc.Register(log, gRPCServer, reservService, bookClient)

	return &App{
		log:        log,
		gRPCServer: gRPCServer,
		port:       port,
	}

}

func (app *App) MustRun() {
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func (app *App) Run() error {
	log := app.log

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", app.port))
	if err != nil {
		return err
	}

	log.Info("Starting gRPC server on port %d", app.port, slog.String("addr", l.Addr().String()))

	if err := app.gRPCServer.Serve(l); err != nil {
		return err
	}

	return nil
}

func (app *App) Stop() {
	log := app.log

	app.gRPCServer.GracefulStop()

	log.Info("Stopping gRPC server on port %d", app.port)
}
