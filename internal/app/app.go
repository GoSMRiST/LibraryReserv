package app

import (
	restapp "ReservationsService/internal/app/rest"
	"ReservationsService/internal/config"
	"ReservationsService/internal/core"
	"context"
	"github.com/GoSMRiST/protosLibary/gen/go/auth"
	"log/slog"
)

type ReservServise interface {
	AddReservation(ctx context.Context, request *core.ReservationRequest) (*core.ReservationResponse, error)
	CloseReservation(ctx context.Context, request *core.ReturnRequest) (*core.ReturnResponse, error)
	CheckReservation(ctx context.Context) (*core.CheckReservResponse, error)
}

type App struct {
	RestServ *restapp.RestReservApp
}

func NewApp(log *slog.Logger,
	cfg *config.Config,
	restService ReservServise,
	authClient auth.AuthClient,
) *App {
	restApp := restapp.NewRestApp(log, cfg, restService, authClient)

	return &App{
		RestServ: restApp,
	}
}
