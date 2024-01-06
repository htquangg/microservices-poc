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
	startBasketEndpoint  endpoint.Endpoint
	cancelBasketEndpoint endpoint.Endpoint
}

func makeBasketEndpoints(c di.Container) basketEndpoints {
	return basketEndpoints{
		startBasketEndpoint:  makeStartBasketEndpoint(c),
		cancelBasketEndpoint: makeCancelBasketEndpoint(c),
	}
}

// start basket
type (
	startBasketRequest struct {
		ID         string `json:"id"`
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
		ID:         uid.GetManager().ID(),
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

		req := request.(startBasketRequest)
		err := app.Commands.StartBasketHandler.Handle(ctx, commands.StartBasket{
			ID:         req.ID,
			CustomerID: req.CustomerID,
		})

		return startBasketResponse{
			ID:  req.ID,
			Err: err,
		}, nil
	}
}

// cancel basket
type (
	cancelBasketRequest struct {
		ID string `json:"customer_id"`
	}
	cancelBasketResponse struct {
		ErrorCode string `json:"error_code,omitempty"`
		Err       error  `json:"err,omitempty"`
	}
)

func decodeCancelBasketRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb_basket.CancelBasketRequest)
	return cancelBasketRequest{
		ID: req.GetId(),
	}, nil
}

func encodeCancelBasketResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(cancelBasketResponse)
	return &pb_basket.CancelBasketResponse{}, resp.Err
}

func makeCancelBasketEndpoint(c di.Container) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		app := di.Get(ctx, constants.ApplicationKey).(*application.Application)

		req := request.(cancelBasketRequest)
		err := app.Commands.CancelBasketHandler.Handle(ctx, commands.CancelBasket{
			ID: req.ID,
		})

		return cancelBasketResponse{
			Err: err,
		}, nil
	}
}
