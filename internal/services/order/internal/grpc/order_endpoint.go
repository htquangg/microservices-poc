package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/order/constants"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/application/commands"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/domain"
	"github.com/htquangg/microservices-poc/internal/services/order/orderpb"
	"github.com/htquangg/microservices-poc/pkg/uid"

	"github.com/go-kit/kit/endpoint"
	"github.com/htquangg/di/v2"
)

type orderEndpoints struct {
	createOrderEndpoint endpoint.Endpoint
	cancelOrderEndpoint endpoint.Endpoint
}

func makeOrderEndpoints(c di.Container) orderEndpoints {
	return orderEndpoints{
		createOrderEndpoint: makeCreateOrderEndpoint(c),
		cancelOrderEndpoint: makeCancelOrderEndpoint(c),
	}
}

// create order
type (
	createOrderRequest struct {
		ID         string         `json:"id"`
		CustomerID string         `json:"customer_id"`
		PaymentID  string         `json:"payment_id"`
		Items      []*domain.Item `json:"items"`
	}
	createOrderResponse struct {
		ID        string `json:"id"`
		ErrorCode string `json:"error_code,omitempty"`
		Err       error  `json:"err,omitempty"`
	}
)

func decodeCreateOrderRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*orderpb.CreateOrderRequest)

	items := make([]*domain.Item, 0, len(req.GetItems()))
	for _, item := range req.GetItems() {
		items = append(
			items,
			domain.NewItem(
				item.GetProductId(),
				item.GetStoreId(),
				item.GetPrice(),
				int(item.GetQuantiy()),
				domain.WithProductName(item.GetProductName()),
				domain.WithStoreName(item.GetStoreName()),
			),
		)
	}

	return createOrderRequest{
		ID:         uid.GetManager().ID(),
		CustomerID: req.GetCustomerId(),
		PaymentID:  req.GetPaymentId(),
		Items:      items,
	}, nil
}

func endcodeCreateOrderResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(createOrderResponse)
	return &orderpb.CreateOrderResponse{
		Id: resp.ID,
	}, resp.Err
}

func makeCreateOrderEndpoint(c di.Container) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		app := di.Get(ctx, constants.ApplicationKey).(*application.Application)

		req := request.(createOrderRequest)
		err := app.Commands.CreateOrderHandler.Handle(ctx, commands.CreateOrder{
			ID:         req.ID,
			CustomerID: req.CustomerID,
			PaymentID:  req.PaymentID,
			Items:      req.Items,
		})

		return createOrderResponse{
			ID:  req.ID,
			Err: err,
		}, nil
	}
}

// cancel order
type (
	cancelOrderRequest struct {
		ID string `json:"customer_id"`
	}
	cancelOrderResponse struct {
		ErrorCode string `json:"error_code,omitempty"`
		Err       error  `json:"err,omitempty"`
	}
)

func decodeCancelOrderRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*orderpb.CancelOrderRequest)
	return cancelOrderRequest{
		ID: req.GetId(),
	}, nil
}

func encodeCancelOrderResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(cancelOrderResponse)
	return &orderpb.CancelOrderResponse{}, resp.Err
}

func makeCancelOrderEndpoint(c di.Container) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		app := di.Get(ctx, constants.ApplicationKey).(*application.Application)

		req := request.(cancelOrderRequest)
		err := app.Commands.CancelOrderHandler.Handle(ctx, commands.CancelOrder{
			ID: req.ID,
		})

		return cancelOrderResponse{
			Err: err,
		}, nil
	}
}
