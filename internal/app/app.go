package app

import (
	grpcapp "ReservationsService/internal/app/grpc"
	"ReservationsService/internal/core"
	"context"
	"log/slog"
)

type App struct {
	GRPCSrv *grpcapp.App
}

type ReservServise interface {
	AddReservation(ctx context.Context, request *core.ReservationRequest) (*core.ReservationResponse, error)
	CloseReservation(ctx context.Context, request *core.ReturnRequest) (*core.ReturnResponse, error)
	CheckReservation(ctx context.Context, userId int) (*core.CheckReservResponse, error)
}

type BookClientService interface {
	CheckBookAvailability(ctx context.Context, req *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error)
}

func NewApp(log *slog.Logger,
	grpcPort int,
	reservService ReservServise,
	bookClient BookClientService,
) *App {
	grpcApp := grpcapp.New(log, grpcPort, reservService, bookClient)

	return &App{
		GRPCSrv: grpcApp,
	}
}
