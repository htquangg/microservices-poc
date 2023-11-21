package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application"
	customerpb "github.com/htquangg/microservices-poc/internal/services/customer/proto"
	"github.com/htquangg/microservices-poc/pkg/database"
	"github.com/htquangg/microservices-poc/pkg/uid"

	grpctransport "github.com/go-kit/kit/transport/grpc"
	"google.golang.org/grpc"
)

var _ customerpb.CustomerServiceServer = (*customerServer)(nil)

type customerServer struct {
	customerpb.UnimplementedCustomerServiceServer
	app *application.Application
	db  *database.DB
	sf  *uid.Sonyflake

	registerCustomer grpctransport.Handler
}

func registerCustomerServer(
	app *application.Application,
	db *database.DB,
	sf *uid.Sonyflake,
	registrar grpc.ServiceRegistrar,
) error {
	s := &customerServer{
		app: app,
		db:  db,
		sf:  sf,
	}

	endpoints := makeCustomerEndpoints(s.app, s.sf)

	s.registerCustomer = grpctransport.NewServer(
		endpoints.registerCustomerEndpoint,
		decodeRegisterCustomerRequest,
		encodeRegisterCustomerResponse,
	)

	customerpb.RegisterCustomerServiceServer(registrar, s)

	return nil
}

func (s *customerServer) RegisterCustomer(
	ctx context.Context,
	request *customerpb.RegisterCustomerRequest,
) (*customerpb.RegisterCustomerResponse, error) {
	var errTx, err error

	var resp interface{}

	errTx = s.db.WithTx(ctx, func(ctx context.Context) error {
		_, resp, err = s.registerCustomer.ServeGRPC(ctx, request)
		return err
	})

	if errTx != nil {
		return nil, errTx
	}

	return resp.(*customerpb.RegisterCustomerResponse), nil
}
