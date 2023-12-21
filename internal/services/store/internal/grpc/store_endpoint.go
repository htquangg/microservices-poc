package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/customer/constants"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/application/commands"
	pb_store "github.com/htquangg/microservices-poc/internal/services/store/proto"
	"github.com/htquangg/microservices-poc/pkg/uid"

	"github.com/go-kit/kit/endpoint"
	"github.com/htquangg/di/v2"
)

type customerEndpoints struct {
	registerCustomerEndpoint endpoint.Endpoint
}

func makeStoreEndpoints(c di.Container) customerEndpoints {
	return customerEndpoints{
		registerCustomerEndpoint: makeAddProductEndpoint(c),
	}
}

type (
	addProductRequest struct {
		StoreID     string  `json:"store_id"`
		Name        string  `json:"name"`
		Description string  `json:"description"`
		SKU         string  `json:"sku"`
		Price       float64 `json:"price"`
	}
	addProductResponse struct {
		ID        string `json:"id"`
		ErrorCode string `json:"error_code,omitempty"`
		Err       error  `json:"err,omitempty"`
	}
)

func decodeAddProductRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb_store.AddProductRequest)
	return addProductRequest{
		StoreID:     req.GetStoreId(),
		Name:        req.GetName(),
		Description: req.GetDescription(),
		SKU:         req.GetSku(),
		Price:       req.GetPrice(),
	}, nil
}

func encodeAddProductResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(addProductResponse)
	return &pb_store.AddProductResponse{
		Id: resp.ID,
	}, resp.Err
}

func makeAddProductEndpoint(c di.Container) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		app := di.Get(ctx, constants.ApplicationKey).(*application.Application)

		id := uid.GetManager().ID()

		req := request.(addProductRequest)
		err := app.Commands.AddProductHandler.Handle(ctx, commands.AddProduct{
			ID:          id,
			StoreID:     req.StoreID,
			Name:        req.Name,
			Description: req.Description,
			SKU:         req.SKU,
			Price:       req.Price,
		})

		return addProductResponse{
			ID:  id,
			Err: err,
		}, nil
	}
}
