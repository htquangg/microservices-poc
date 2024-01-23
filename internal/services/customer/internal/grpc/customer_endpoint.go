package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/customer/constants"
	"github.com/htquangg/microservices-poc/internal/services/customer/customerpb"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application/commands"
	"github.com/htquangg/microservices-poc/pkg/uid"

	"github.com/go-kit/kit/endpoint"
	"github.com/htquangg/di/v2"
)

type customerEndpoints struct {
	registerCustomerEndpoint endpoint.Endpoint
}

func makeCustomerEndpoints(ctn di.Container) customerEndpoints {
	return customerEndpoints{
		registerCustomerEndpoint: makeRegisterCustomerEndpoint(ctn),
	}
}

type (
	registerCustomerRequest struct {
		ID    string `json:"id"`
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
	req := grpcReq.(*customerpb.RegisterCustomerRequest)
	return registerCustomerRequest{
		ID:    uid.GetManager().ID(),
		Name:  req.GetName(),
		Phone: req.GetPhone(),
		Email: req.GetEmail(),
	}, nil
}

func encodeRegisterCustomerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(registerCustomerResponse)
	return &customerpb.RegisterCustomerResponse{
		Id: resp.ID,
	}, resp.Err
}

func makeRegisterCustomerEndpoint(ctn di.Container) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		app := ctn.Get(constants.ApplicationKey).(*application.Application)

		req := request.(registerCustomerRequest)
		err := app.Commands.RegisterCustomerHandler.Handle(ctx, commands.RegisterCustomer{
			ID:    req.ID,
			Name:  req.Name,
			Phone: req.Phone,
			Email: req.Email,
		})

		return registerCustomerResponse{
			ID:  req.ID,
			Err: err,
		}, nil
	}
}
