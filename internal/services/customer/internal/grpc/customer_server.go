package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/customer/customerpb"
	"github.com/htquangg/microservices-poc/pkg/database"

	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"github.com/htquangg/di/v2"
	"google.golang.org/grpc"
)

var _ customerpb.CustomerServiceServer = (*customerServer)(nil)

type customerServer struct {
	customerpb.UnimplementedCustomerServiceServer

	ctn di.Container
	db  database.DB

	registerCustomer grpc_transport.Handler
}

func registerCustomerServer(
	ctn di.Container,
	db database.DB,
	registrar grpc.ServiceRegistrar,
) error {
	endpoints := makeCustomerEndpoints(ctn)

	customerpb.RegisterCustomerServiceServer(registrar, customerServer{
		ctn: ctn,
		db:  db,
		registerCustomer: grpc_transport.NewServer(
			endpoints.registerCustomerEndpoint,
			decodeRegisterCustomerRequest,
			encodeRegisterCustomerResponse,
		),
	})

	return nil
}

func (s customerServer) RegisterCustomer(
	ctx context.Context,
	request *customerpb.RegisterCustomerRequest,
) (*customerpb.RegisterCustomerResponse, error) {
	var resp interface{}

	err := s.db.WithTx(ctx, func(ctx context.Context) (err error) {
		_, resp, err = s.registerCustomer.ServeGRPC(ctx, request)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp.(*customerpb.RegisterCustomerResponse), nil
}
