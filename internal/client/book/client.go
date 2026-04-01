package bookclient

import (
	"ReservationsService/internal/core"
	"context"
	"fmt"
	bookclient "github.com/GoSMRiST/protosLibary/gen/go/book"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type BookClient struct {
	client bookclient.BookClient
}

func NewBookClient(client bookclient.BookClient) *BookClient {
	return &BookClient{client: client}
}

func (c *BookClient) CheckAvailability(ctx context.Context, req *core.CheckAvailabilityRequest) (*core.CheckAvailabilityResponse, error) {
	pbReq := &bookclient.CheckRequest{
		Author: req.Author,
		Title:  req.Title,
	}

	response := &core.CheckAvailabilityResponse{
		Result: false,
	}

	result, err := c.client.CheckAvailability(ctx, pbReq)
	if err != nil {
		st, ok := status.FromError(err)
		if ok {
			switch st.Code() {

			case codes.InvalidArgument:
				return response, core.ErrInvalidInput

			case codes.Internal:
				return response, fmt.Errorf("book service internal error")

			case codes.NotFound:
				return response, core.ErrNotFound

			default:
				return response, err
			}
		}

		return response, err
	}

	response.Result = result.Result
	return response, nil
}
