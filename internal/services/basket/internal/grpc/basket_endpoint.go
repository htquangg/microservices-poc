package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/basket/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/application/commands"
	pb_basket "github.com/htquangg/microservices-poc/internal/services/basket/proto"
	"github.com/htquangg/microservices-poc/internal/services/customer/constants"
	"github.com/htquangg/microservices-poc/pkg/uid"

	"github.com/go-kit/kit/endpoint"
	"github.com/htquangg/di/v2"
)

type basketEndpoints struct {
	startBasketEndpoint endpoint.Endpoint
}

func makeBasketEndpoints(c di.Container) basketEndpoints {
	return basketEndpoints{
		startBasketEndpoint: makeStartBasketEndpoint(c),
	}
}

type (
	startBasketRequest struct {
		CustomerID string `json:"customer_id"`
	}
	startBasketResponse struct {
		ID        string `json:"id"`
		ErrorCode string `json:"error_code,omitempty"`
		Err       error  `json:"err,omitempty"`
	}
)

func decodeStartBasketRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb_basket.StartBasketRequest)
	return startBasketRequest{
		CustomerID: req.GetCustomerId(),
	}, nil
}

func encodeStartBasketResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(startBasketResponse)
	return &pb_basket.StartBasketResponse{
		Id: resp.ID,
	}, resp.Err
}

func makeStartBasketEndpoint(c di.Container) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		app := di.Get(ctx, constants.ApplicationKey).(*application.Application)

		id := uid.GetManager().ID()

		req := request.(startBasketRequest)
		err := app.Commands.StartBasketHandler.Handle(ctx, commands.StartBasket{
			ID:         id,
			CustomerID: req.CustomerID,
		})

		return startBasketResponse{
			ID:  id,
			Err: err,
		}, nil
	}
}
