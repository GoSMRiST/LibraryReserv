package rest

import (
	"ReservationsService/internal/config"
	"ReservationsService/internal/core"
	"ReservationsService/internal/middleware"
	"ReservationsService/internal/transport/rest"
	"context"
	"errors"
	"github.com/GoSMRiST/protosLibary/gen/go/auth"
	"github.com/gin-gonic/gin"
	"log/slog"
	"net/http"
)

type RestReservServise interface {
	AddReservation(ctx context.Context, request *core.ReservationRequest) (*core.ReservationResponse, error)
	CloseReservation(ctx context.Context, request *core.ReturnRequest) (*core.ReturnResponse, error)
	CheckReservation(ctx context.Context, userId int) (*core.CheckReservResponse, error)
}

type RestReservApp struct {
	log        *slog.Logger
	restServer *http.Server
	hostAddr   string
}

func NewRestApp(log *slog.Logger, cfg *config.Config, service RestReservServise, grpcAuthClient auth.AuthClient) *RestReservApp {
	reservHandler := rest.NewRestReservServer(log, service)

	engine := gin.Default()
	engine.Use(middleware.AuthMiddleware(grpcAuthClient))

	reservHandler.RegisterRoutes(engine)

	srv := &http.Server{
		Addr:         cfg.HostAddress,
		Handler:      engine,
		ReadTimeout:  cfg.ServTimeout,
		WriteTimeout: cfg.ServTimeout,
		IdleTimeout:  cfg.ServTimeout,
	}

	return &RestReservApp{
		log:        log,
		restServer: srv,
		hostAddr:   cfg.HostAddress,
	}
}

func (a *RestReservApp) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

func (a *RestReservApp) Run() error {
	a.log.Info("starting http server on ",
		"addr:", a.hostAddr,
	)

	if err := a.restServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.log.Error("REST server failed", "error", err)
		return err
	}

	return nil
}

func (a *RestReservApp) Stop(ctx context.Context) error {
	if err := a.restServer.Shutdown(ctx); err != nil {
		a.log.Error("REST shutdown error", "error", err)
		return err
	}

	a.log.Info("REST server stopped")

	return nil
}
