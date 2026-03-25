package reserv

import (
	"ReservationsService/internal/core"
	"context"
	reserv "github.com/GoSMRiST/protosLibary/gen/go/reserv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	"time"
)

type ReservServise interface {
	AddReservation(ctx context.Context, request *core.ReservationRequest) (*core.ReservationResponse, error)
	CloseReservation(ctx context.Context, request *core.ReturnRequest) (*core.ReturnResponse, error)
	CheckReservation(ctx context.Context, userId int) (*core.CheckReservResponse, error)
}

type BookClientService interface {
	CheckBookAvailability(ctx context.Context, req *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error)
}

type ServerAPI struct {
	logger *slog.Logger
	reserv.UnimplementedReservServer
	service    ReservServise
	bookClient BookClientService
}

func Register(logger *slog.Logger, gRPC *grpc.Server, service ReservServise, bookClient BookClientService) {
	reserv.RegisterReservServer(gRPC, &ServerAPI{logger: logger, service: service, bookClient: bookClient})
}

func (s *ServerAPI) ReservBook(ctx context.Context, req *reserv.ReservRequest) (*reserv.ReservResponse, error) {
	coreRequest := core.ReservationRequest{
		UserID: int(req.GetUserId()),
		Author: req.GetAuthor(),
		Title:  req.GetTitle(),
	}

	requestForBookClient := &core.CheckAvailabilityRequest{
		Author: req.GetAuthor(),
		Title:  req.GetTitle(),
	}

	bookAvailability, err := s.bookClient.CheckBookAvailability(ctx, requestForBookClient)
	if err != nil {
		return &reserv.ReservResponse{
			Result: bookAvailability.Result,
		}, status.Error(codes.Internal, err.Error())
	}

	if bookAvailability.Result != true {
		return &reserv.ReservResponse{
			Result: bookAvailability.Result,
		}, status.Error(codes.NotFound, "Book is not found")
	}

	resp, err := s.service.AddReservation(ctx, &coreRequest)
	if err != nil {

		st, ok := status.FromError(err)
		if ok {
			return &reserv.ReservResponse{
				Result: resp.ReservationStatus,
			}, st.Err()
		}

		return &reserv.ReservResponse{
			Result: resp.ReservationStatus,
		}, status.Error(codes.Internal, err.Error())
	}

	s.logger.Info("Операция резевации успешно выполнена",
		"user_id", req.GetUserId(),
	)

	return &reserv.ReservResponse{
		Result: resp.ReservationStatus,
	}, nil
}

func (s *ServerAPI) ReturnBook(ctx context.Context, req *reserv.ReturnRequest) (*reserv.ReturnResponse, error) {
	coreRequest := core.ReturnRequest{
		UserID: int(req.GetUserId()),
		Author: req.GetAuthor(),
		Title:  req.GetTitle(),
	}

	resp, err := s.service.CloseReservation(ctx, &coreRequest)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return &reserv.ReturnResponse{
				Result: resp.Status,
			}, st.Err()
		}

		return &reserv.ReturnResponse{
			Result: resp.Status,
		}, status.Error(codes.Internal, err.Error())
	}

	s.logger.Info("Пользователь успешно сдал книгу ",
		"user_id", req.GetUserId(),
		"author", req.GetAuthor(),
		"title", req.GetTitle(),
	)

	return &reserv.ReturnResponse{
		Result: resp.Status,
	}, nil
}

func (s *ServerAPI) CheckReservation(ctx context.Context, req *reserv.CheckReservRequest) (*reserv.CheckReservResponse, error) {
	resp, err := s.service.CheckReservation(ctx, int(req.GetUserID()))
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			return nil, st.Err()
		}

		return nil, status.Error(codes.Internal, err.Error())
	}

	pbReservations := make([]*reserv.Reservation, 0, len(resp.Reservations))

	for _, r := range resp.Reservations {
		var returnAt string
		if r.ReturnAt != nil {
			returnAt = r.ReturnAt.Format(time.RFC3339)
		}

		pbReservations = append(pbReservations, &reserv.Reservation{
			ReservationId: int64(r.ReservationID),
			UserId:        int64(r.UserID),
			Author:        r.Author,
			Title:         r.Title,
			TakenAt:       r.TakenAt.String(),
			ReturnAt:      returnAt,
		})
	}

	s.logger.Info("user reservations found",
		"user_id", req.GetUserID(),
		"count", len(pbReservations),
	)

	return &reserv.CheckReservResponse{
		Reservations: pbReservations,
	}, nil
}
