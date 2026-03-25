package main

import (
	"ReservationsService/internal/app"
	bookclient "ReservationsService/internal/client/book"
	"ReservationsService/internal/config"
	"ReservationsService/internal/repository"
	"ReservationsService/internal/services/bookclientserv"
	"ReservationsService/internal/services/reservserv"
	"context"
	"fmt"
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
		log.Info("failed to connect to database")
		panic(err)
	}

	err = db.CreateTable(ctx)
	if err != nil {
		log.Info("failed to create table")
		panic(err)
	}

	grpcBook, err := grpc.Dial("localhost:44045", grpc.WithInsecure())
	if err != nil {
		log.Info("failed to connect to book service")
		panic(err)
	}

	grpcBookClient := book.NewBookClient(grpcBook)

	bookClient := bookclient.NewBookClient(grpcBookClient)
	if bookClient == nil {
		log.Info("failed to create book client")
		panic(err)
	}

	reservService := reservserv.NewReservService(log, db)
	bookClientService := bookclientserv.NewBookClientServ(bookClient)

	application := app.NewApp(log, cfg.ServPort, reservService, bookClientService)

	go application.GRPCSrv.MustRun()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	stopSign := <-stop

	log.Info("stop signal", slog.Any("signal", stopSign))

	application.GRPCSrv.Stop()

	log.Info("application stopped")

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
