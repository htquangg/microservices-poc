package grpc

import (
	"context"

	pb_customer "github.com/htquangg/microservices-poc/internal/services/customer/proto"
	"github.com/htquangg/microservices-poc/pkg/database"

	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"github.com/htquangg/di/v2"
	"google.golang.org/grpc"
)

var _ pb_customer.CustomerServiceServer = (*customerServer)(nil)

type customerServer struct {
	pb_customer.UnimplementedCustomerServiceServer

	c  di.Container
	db database.DB

	registerCustomer grpc_transport.Handler
}

func registerCustomerServer(
	c di.Container,
	db database.DB,
	registrar grpc.ServiceRegistrar,
) error {
	endpoints := makeCustomerEndpoints(c)

	pb_customer.RegisterCustomerServiceServer(registrar, customerServer{
		c:  c,
		db: db,
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
	request *pb_customer.RegisterCustomerRequest,
) (*pb_customer.RegisterCustomerResponse, error) {
	ctx = s.c.Scoped(ctx)

	var resp interface{}

	err := s.db.WithTx(ctx, func(ctx context.Context) (err error) {
		_, resp, err = s.registerCustomer.ServeGRPC(ctx, request)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp.(*pb_customer.RegisterCustomerResponse), nil
}
