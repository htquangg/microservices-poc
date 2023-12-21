package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/customer/constants"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application/commands"
	pb_customer "github.com/htquangg/microservices-poc/internal/services/customer/proto"
	"github.com/htquangg/microservices-poc/pkg/uid"

	"github.com/go-kit/kit/endpoint"
	"github.com/htquangg/di/v2"
)

type customerEndpoints struct {
	registerCustomerEndpoint endpoint.Endpoint
}

func makeCustomerEndpoints(c di.Container) customerEndpoints {
	return customerEndpoints{
		registerCustomerEndpoint: makeRegisterCustomerEndpoint(c),
	}
}

type (
	registerCustomerRequest struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
		Email string `json:"email"`
	}
	registerCustomerResponse struct {
		ID        string `json:"id"`
		ErrorCode string `json:"error_code,omitempty"`
		Err       error  `json:"err,omitempty"`
	}
)

func decodeRegisterCustomerRequest(_ context.Context, grpcReq interface{}) (interface{}, error) {
	req := grpcReq.(*pb_customer.RegisterCustomerRequest)
	return registerCustomerRequest{
		Name:  req.GetName(),
		Phone: req.GetPhone(),
		Email: req.GetEmail(),
	}, nil
}

func encodeRegisterCustomerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(registerCustomerResponse)
	return &pb_customer.RegisterCustomerResponse{
		Id: resp.ID,
	}, resp.Err
}

func makeRegisterCustomerEndpoint(c di.Container) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		app := di.Get(ctx, constants.ApplicationKey).(*application.Application)

		id := uid.GetManager().ID()

		req := request.(registerCustomerRequest)
		err := app.Commands.RegisterCustomerHandler.Handle(ctx, commands.RegisterCustomer{
			ID:    id,
			Name:  req.Name,
			Phone: req.Phone,
			Email: req.Email,
		})

		return registerCustomerResponse{
			ID:  id,
			Err: err,
		}, nil
	}
}
