package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/customer/constants"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application/command"
	customerpb "github.com/htquangg/microservices-poc/internal/services/customer/proto"
	"github.com/htquangg/microservices-poc/pkg/uid"

	"github.com/go-kit/kit/endpoint"
	"github.com/htquangg/di/v2"
)

type customerEndpoints struct {
	registerCustomerEndpoint endpoint.Endpoint
}

func makeCustomerEndpoints(c di.Container, sf *uid.Sonyflake) customerEndpoints {
	return customerEndpoints{
		registerCustomerEndpoint: makeRegisterCustomerEndpoint(c, sf),
	}
}

type (
	registerCustomerRequest struct {
		Name  string `json:"name"`
		Phone string `json:"phone"`
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
		Name:  req.GetName(),
		Phone: req.GetPhone(),
	}, nil
}

func encodeRegisterCustomerResponse(_ context.Context, response interface{}) (interface{}, error) {
	resp := response.(registerCustomerResponse)
	return &customerpb.RegisterCustomerResponse{
		Id: resp.ID,
	}, resp.Err
}

func makeRegisterCustomerEndpoint(c di.Container, sf *uid.Sonyflake) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		app := di.Get(ctx, constants.ApplicationKey).(*application.Application)

		req := request.(registerCustomerRequest)

		customerID := sf.ID()

		err := app.Commands.RegisterCustomerHandler.Handle(ctx, command.RegisterCustomer{
			ID:    customerID,
			Name:  req.Name,
			Phone: req.Phone,
		})

		return registerCustomerResponse{
			ID:  customerID,
			Err: err,
		}, nil
	}
}
