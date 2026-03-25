package reservserv

import (
	"ReservationsService/internal/core"
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
)

type ReservService struct {
	repo DBRepository
}

func NewReservService(log *slog.Logger, repo DBRepository) *ReservService {
	return &ReservService{repo}
}

type DBRepository interface {
	AddReservation(ctx context.Context, request *core.ReservationRequest) (*core.ReservationResponse, error)
	CloseReservation(ctx context.Context, request *core.ReturnRequest) (*core.ReturnResponse, error)
	CheckReservation(ctx context.Context, userId int) (*core.CheckReservResponse, error)
}

func (sr *ReservService) AddReservation(ctx context.Context, request *core.ReservationRequest) (*core.ReservationResponse, error) {
	if request.Author == "" || request.Title == "" || request.UserID <= 0 {
		return &core.ReservationResponse{}, status.Error(codes.InvalidArgument, "Author, Title, or UserId is empty")
	}

	return sr.repo.AddReservation(ctx, request)
}

func (sr *ReservService) CloseReservation(ctx context.Context, request *core.ReturnRequest) (*core.ReturnResponse, error) {
	if request.Author == "" || request.Title == "" || request.UserID <= 0 {
		return &core.ReturnResponse{}, status.Error(codes.InvalidArgument, "Author, Title, or UserId is empty")
	}

	return sr.repo.CloseReservation(ctx, request)
}

func (sr *ReservService) CheckReservation(ctx context.Context, userId int) (*core.CheckReservResponse, error) {
	if userId <= 0 {
		return &core.CheckReservResponse{}, status.Error(codes.InvalidArgument, "UserID is empty")
	}

	return sr.repo.CheckReservation(ctx, userId)
}
