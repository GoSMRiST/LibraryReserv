package main

import (
	"ReservationsService/internal/app"
	bookclient "ReservationsService/internal/client/book"
	"ReservationsService/internal/config"
	"ReservationsService/internal/repository"
	"ReservationsService/internal/services/bookclientserv"
	"ReservationsService/internal/services/restserv"
	"context"
	"fmt"
	"github.com/GoSMRiST/protosLibary/gen/go/auth"
	"github.com/GoSMRiST/protosLibary/gen/go/book"
	"google.golang.org/grpc"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	ctx := context.Background()

	cfg := config.InitConfig()

	log := setupLogger(cfg.LogLevel)

	log.Info("starting", slog.Any("config", cfg))

	connString := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	db, err := repository.InitDataBase(ctx, connString)
	if err != nil {
		log.Error("failed to connect to database")
		panic(err)
	}

	err = db.CreateTable(ctx)
	if err != nil {
		log.Error("failed to create table")
		panic(err)
	}

	grpcBook, err := grpc.Dial("localhost:44045", grpc.WithInsecure())
	if err != nil {
		log.Error("failed to connect to book service")
		panic(err)
	}

	grpcBookClient := book.NewBookClient(grpcBook)

	bookClient := bookclient.NewBookClient(grpcBookClient)
	if bookClient == nil {
		log.Error("failed to create book client")
		panic(err)
	}
	defer func() {
		if err := grpcBook.Close(); err != nil {
			log.Error("failed to close grpc book service", "error", err)
		}
	}()

	grpcAuth, err := grpc.Dial("localhost:44046", grpc.WithInsecure())
	if err != nil {
		log.Error("failed to connect to auth service", "error", err)
		panic(err)
	}
	defer func() {
		if err := grpcAuth.Close(); err != nil {
			log.Info("failed to close grpc auth service: ", err)
		}
	}()

	bookClientService := bookclientserv.NewBookClientServ(bookClient)
	grpcAuthClient := auth.NewAuthClient(grpcAuth)

	reservRestService := restserv.NewRestReservService(log, db, bookClientService)

	application := app.NewApp(log, cfg, reservRestService, grpcAuthClient)

	go application.RestServ.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	sign := <-stop

	log.Info("received signal", sign.String())

	ctxTimeout, cancel := context.WithTimeout(ctx, cfg.ServTimeout)
	defer cancel()

	if err := application.RestServ.Stop(ctxTimeout); err != nil {
		log.Error("failed to stop REST server", "error", err)
		return
	}

	log.Info("rest server stopped")

	if db.Close(ctx) != nil {
		log.Info("fail to close database")
	}
	log.Info("database is closed")
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envDev:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)

	case envProd:
		logger = slog.New(
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)

	default:
		logger = slog.New(
			slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		)
	}

	return logger
}
