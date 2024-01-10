package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/basket/basketpb"
	"github.com/htquangg/microservices-poc/pkg/database"

	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"github.com/htquangg/di/v2"
	"google.golang.org/grpc"
)

var _ basketpb.BasketServiceServer = (*basketServer)(nil)

type basketServer struct {
	basketpb.UnimplementedBasketServiceServer

	c  di.Container
	db database.DB

	startBasket  grpc_transport.Handler
	cancelBasket grpc_transport.Handler
}

func registerBasketServer(
	c di.Container,
	db database.DB,
	registrar grpc.ServiceRegistrar,
) error {
	endpoints := makeBasketEndpoints(c)

	basketpb.RegisterBasketServiceServer(registrar, basketServer{
		c:  c,
		db: db,
		startBasket: grpc_transport.NewServer(
			endpoints.startBasketEndpoint,
			decodeStartBasketRequest,
			encodeStartBasketResponse,
		),
		cancelBasket: grpc_transport.NewServer(
			endpoints.cancelBasketEndpoint,
			decodeCancelBasketRequest,
			encodeCancelBasketResponse,
		),
	})

	return nil
}

func (s basketServer) StartBasket(
	ctx context.Context,
	request *basketpb.StartBasketRequest,
) (*basketpb.StartBasketResponse, error) {
	ctx = s.c.Scoped(ctx)

	var resp interface{}

	err := s.db.WithTx(ctx, func(ctx context.Context) (err error) {
		_, resp, err = s.startBasket.ServeGRPC(ctx, request)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp.(*basketpb.StartBasketResponse), nil
}

func (s basketServer) CancelBasket(
	ctx context.Context,
	request *basketpb.CancelBasketRequest,
) (*basketpb.CancelBasketResponse, error) {
	ctx = s.c.Scoped(ctx)

	var resp interface{}

	err := s.db.WithTx(ctx, func(ctx context.Context) (err error) {
		_, resp, err = s.cancelBasket.ServeGRPC(ctx, request)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp.(*basketpb.CancelBasketResponse), nil
}
