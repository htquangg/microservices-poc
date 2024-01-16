package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/basket/basketpb"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/application/commands"
	"github.com/htquangg/microservices-poc/internal/services/customer/constants"
	"github.com/htquangg/microservices-poc/pkg/uid"

	"github.com/go-kit/kit/endpoint"
	"github.com/htquangg/di/v2"
)

type basketEndpoints struct {
	startBasketEndpoint    endpoint.Endpoint
	cancelBasketEndpoint   endpoint.Endpoint
	checkoutBasketEndpoint endpoint.Endpoint
	addItemEndpoint        endpoint.Endpoint
	removeItemEndpoint     endpoint.Endpoint
}

func makeBasketEndpoints(c di.Container) basketEndpoints {
	return basketEndpoints{
		startBasketEndpoint:    makeStartBasketEndpoint(c),
		cancelBasketEndpoint:   makeCancelBasketEndpoint(c),
		checkoutBasketEndpoint: makeCheckoutBasketEndpoint(c),
		addItemEndpoint:        makeAddItemEndpoint(c),
		removeItemEndpoint:     makeRemoveItemEndpoint(c),
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
	req := grpcReq.(*basketpb.StartBasketRequest)
	return startBasketRequest{
		ID:         uid.GetManager().ID(),
		CustomerID: req.GetCustomerId(),
	}, nil
}

func encodeStartBasketResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(startBasketResponse)
	return &basketpb.StartBasketResponse{
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
		ID string `json:"id"`
	}
	cancelBasketResponse struct {
		ErrorCode string `json:"error_code,omitempty"`
		Err       error  `json:"err,omitempty"`
	}
)

func decodeCancelBasketRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*basketpb.CancelBasketRequest)
	return cancelBasketRequest{
		ID: req.GetId(),
	}, nil
}

func encodeCancelBasketResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(cancelBasketResponse)
	return &basketpb.CancelBasketResponse{}, resp.Err
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

// checkout basket
type (
	checkoutBasketRequest struct {
		ID        string `json:"customer_id"`
		PaymentID string `json:"payment_id"`
	}
	checkoutBasketResponse struct {
		ErrorCode string `json:"error_code,omitempty"`
		Err       error  `json:"err,omitempty"`
	}
)

func decodeCheckoutBasketRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*basketpb.CheckoutBasketRequest)
	return checkoutBasketRequest{
		ID:        req.GetId(),
		PaymentID: req.GetPaymentId(),
	}, nil
}

func encodeCheckoutBasketResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(checkoutBasketResponse)
	return &basketpb.CheckoutBasketResponse{}, resp.Err
}

func makeCheckoutBasketEndpoint(c di.Container) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		app := di.Get(ctx, constants.ApplicationKey).(*application.Application)

		req := request.(checkoutBasketRequest)
		err := app.Commands.CheckoutBasketHandler.Handle(ctx, commands.CheckoutBasket{
			ID:        req.ID,
			PaymentID: req.PaymentID,
		})

		return checkoutBasketResponse{
			Err: err,
		}, nil
	}
}

// add item to basket
type (
	addItemRequest struct {
		ID        string `json:"id"`
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	addItemResponse struct {
		ErrorCode string `json:"error_code,omitempty"`
		Err       error  `json:"err,omitempty"`
	}
)

func decodeAddItemRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*basketpb.AddItemRequest)
	return addItemRequest{
		ID:        req.GetId(),
		ProductID: req.GetProductId(),
		Quantity:  int(req.GetQuantity()),
	}, nil
}

func encodeAddItemResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(addItemResponse)
	return &basketpb.AddItemResponse{}, resp.Err
}

func makeAddItemEndpoint(c di.Container) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		app := di.Get(ctx, constants.ApplicationKey).(*application.Application)

		req := request.(addItemRequest)
		err := app.Commands.AddItemHandler.Handle(ctx, commands.AddItem{
			ID:        req.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		})

		return addItemResponse{
			Err: err,
		}, nil
	}
}

// remote item from basket
type (
	removeItemRequest struct {
		ID        string `json:"id"`
		ProductID string `json:"product_id"`
		Quantity  int    `json:"quantity"`
	}
	removeItemResponse struct {
		ErrorCode string `json:"error_code,omitempty"`
		Err       error  `json:"err,omitempty"`
	}
)

func decodeRemoveItemRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*basketpb.RemoveItemRequest)
	return removeItemRequest{
		ID:        req.GetId(),
		ProductID: req.GetProductId(),
		Quantity:  int(req.GetQuantity()),
	}, nil
}

func encodeRemoveItemResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(removeItemResponse)
	return &basketpb.RemoveItemResponse{}, resp.Err
}

func makeRemoveItemEndpoint(c di.Container) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		app := di.Get(ctx, constants.ApplicationKey).(*application.Application)

		req := request.(removeItemRequest)
		err := app.Commands.RemoveItemHandler.Handle(ctx, commands.RemoveItem{
			ID:        req.ID,
			ProductID: req.ProductID,
			Quantity:  req.Quantity,
		})

		return removeItemResponse{
			Err: err,
		}, nil
	}
}
