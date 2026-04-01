package bookclientserv

import (
	"ReservationsService/internal/core"
	"context"
	"errors"
)

type BookClientServ struct {
	client ClientInterface
}

type ClientInterface interface {
	CheckAvailability(ctx context.Context, req *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error)
}

func NewBookClientServ(client ClientInterface) *BookClientServ {
	return &BookClientServ{client: client}
}

func (b *BookClientServ) CheckBookAvailability(ctx context.Context, req *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error) {
	resp, err := b.client.CheckAvailability(ctx, req)
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
	return resp, nil
}
