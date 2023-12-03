package grpc

import (
	"context"

	customerpb "github.com/htquangg/microservices-poc/internal/services/customer/proto"
	"github.com/htquangg/microservices-poc/pkg/database"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"github.com/htquangg/di/v2"
	"google.golang.org/grpc"
)

var _ customerpb.CustomerServiceServer = (*customerServer)(nil)

type customerServer struct {
	customerpb.UnimplementedCustomerServiceServer

	c  di.Container
	db *database.DB

	registerCustomer grpctransport.Handler
}

func registerCustomerServer(
	c di.Container,
	db *database.DB,
	registrar grpc.ServiceRegistrar,
) error {
	endpoints := makeCustomerEndpoints(c)

	customerpb.RegisterCustomerServiceServer(registrar, customerServer{
		c:  c,
		db: db,
		registerCustomer: grpctransport.NewServer(
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
	ctx = s.c.Scoped(ctx)

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
