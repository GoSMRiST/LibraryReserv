package restserv

import (
	"ReservationsService/internal/core"
	"context"
	"errors"
	"log/slog"
)

type DBInterface interface {
	AddReservation(ctx context.Context, request *core.ReservationRequest) (*core.ReservationResponse, error)
	CloseReservation(ctx context.Context, request *core.ReturnRequest) (*core.ReturnResponse, error)
	CheckReservation(ctx context.Context, userId int) (*core.CheckReservResponse, error)
}

type BookClientService interface {
	CheckBookAvailability(ctx context.Context, req *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error)
}

type RestReservService struct {
	log        *slog.Logger
	db         DBInterface
	bookClient BookClientService
}

func NewRestReservService(log *slog.Logger, db DBInterface, bookClient BookClientService) *RestReservService {
	return &RestReservService{
		log:        log,
		db:         db,
		bookClient: bookClient,
	}
}

func (s *RestReservService) AddReservation(ctx context.Context, request *core.ReservationRequest) (*core.ReservationResponse, error) {
	tokenData, ok := ctx.Value(core.TokenDataKey).(core.TokenData)
	if ok == false {
		return nil, core.ErrUnauthorized
	}

	request.UserID = tokenData.UserID

	if tokenData.Role == "user" {
		return nil, core.ErrNoRights
	}

	requestForBookClient := &core.CheckAvailabilityRequest{
		Author: request.Author,
		Title:  request.Title,
	}

	bookAvailability, err := s.bookClient.CheckBookAvailability(ctx, requestForBookClient)
	if err != nil {
		switch {
		case errors.Is(err, core.ErrNotFound):
			return nil, core.ErrNotFound

		case errors.Is(err, core.ErrInvalidInput):
			return nil, core.ErrInvalidInput

		default:
			return nil, err
		}
	}

	if bookAvailability.Result != true {
		return nil, core.ErrBookNotFound
	}

	s.log.Info(requestForBookClient.Author, requestForBookClient.Title)

	if request.Author == "" || request.Title == "" || request.UserID <= 0 {
		return nil, core.ErrInvalidInput
	}

	return s.db.AddReservation(ctx, request)
}

func (s *RestReservService) CloseReservation(ctx context.Context, request *core.ReturnRequest) (*core.ReturnResponse, error) {
	tokenData, ok := ctx.Value(core.TokenDataKey).(core.TokenData)
	if ok == false {
		return nil, core.ErrUnauthorized
	}

	if tokenData.Role == "user" {
		return nil, core.ErrNoRights
	}

	if request.ReservationID <= 0 {
		return nil, core.ErrInvalidInput
	}

	return s.db.CloseReservation(ctx, request)
}

func (s *RestReservService) CheckReservation(ctx context.Context) (*core.CheckReservResponse, error) {
	tokenData, ok := ctx.Value(core.TokenDataKey).(core.TokenData)
	if !ok {
		s.log.Warn("token data not found in context")
		return nil, core.ErrUnauthorized
	}

	s.log.Info("TokenData:", tokenData)

	if tokenData.UserID < 1 {
		return nil, core.ErrInvalidInput
	}

	return s.db.CheckReservation(ctx, tokenData.UserID)
}
