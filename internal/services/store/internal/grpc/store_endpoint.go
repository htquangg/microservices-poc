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

type storeEndpoints struct {
	createStoreEndpoint      endpoint.Endpoint
	rebrandStoreEndpoint     endpoint.Endpoint
	registerCustomerEndpoint endpoint.Endpoint
}

func makeStoreEndpoints(c di.Container) storeEndpoints {
	return storeEndpoints{
		createStoreEndpoint:      makeCreateStoreEndpoint(c),
		rebrandStoreEndpoint:     makeRebrandStoreEndpoint(c),
		registerCustomerEndpoint: makeAddProductEndpoint(c),
	}
}

// create store
type (
	createStoreRequest struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	createStoreResponse struct {
		ID        string `json:"id"`
		ErrorCode string `json:"error_code,omitempty"`
		Err       error  `json:"err,omitempty"`
	}
)

func decodeCreateStoreRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb_store.CreateStoreRequest)
	return createStoreRequest{
		ID:   uid.GetManager().ID(),
		Name: req.GetName(),
	}, nil
}

func encodeCreateStoreResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(createStoreResponse)
	return &pb_store.CreateStoreResponse{
		Id: resp.ID,
	}, resp.Err
}

func makeCreateStoreEndpoint(c di.Container) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		app := di.Get(ctx, constants.ApplicationKey).(*application.Application)

		req := request.(createStoreRequest)
		err := app.Commands.CreateStoreHandler.Handle(ctx, commands.CreateStore{
			ID:   req.ID,
			Name: req.Name,
		})

		return createStoreResponse{
			ID:  req.ID,
			Err: err,
		}, nil
	}
}

// rebrand store
type (
	rebrandStoreRequest struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	rebrandStoreResponse struct {
		ID        string `json:"id"`
		ErrorCode string `json:"error_code,omitempty"`
		Err       error  `json:"err,omitempty"`
	}
)

func decodeRebrandStoreRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb_store.RebrandStoreRequest)
	return rebrandStoreRequest{
		ID:   req.Id,
		Name: req.GetName(),
	}, nil
}

func encodeRebrandStoreResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(rebrandStoreResponse)
	return &pb_store.RebrandStoreResponse{}, resp.Err
}

func makeRebrandStoreEndpoint(c di.Container) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		app := di.Get(ctx, constants.ApplicationKey).(*application.Application)

		req := request.(rebrandStoreRequest)
		err := app.Commands.RebrandStoreHandler.Handle(ctx, commands.RebrandStore{
			ID:   req.ID,
			Name: req.Name,
		})

		return rebrandStoreResponse{
			Err: err,
		}, nil
	}
}

// add product
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
