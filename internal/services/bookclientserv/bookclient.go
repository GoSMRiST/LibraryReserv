package bookclientserv

import (
	"ReservationsService/internal/core"
	"context"
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
	return b.client.CheckAvailability(ctx, req)
}
