package grpc

import (
	"context"

	pb_basket "github.com/htquangg/microservices-poc/internal/services/basket/proto"
	"github.com/htquangg/microservices-poc/pkg/database"

	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"github.com/htquangg/di/v2"
	"google.golang.org/grpc"
)

var _ pb_basket.BasketServiceServer = (*customerServer)(nil)

type customerServer struct {
	pb_basket.UnimplementedBasketServiceServer

	c  di.Container
	db database.DB

	startBasket grpc_transport.Handler
}

func registerBasketServer(
	c di.Container,
	db database.DB,
	registrar grpc.ServiceRegistrar,
) error {
	endpoints := makeBasketEndpoints(c)

	pb_basket.RegisterBasketServiceServer(registrar, customerServer{
		c:  c,
		db: db,
		startBasket: grpc_transport.NewServer(
			endpoints.startBasketEndpoint,
			decodeStartBasketRequest,
			encodeStartBasketResponse,
		),
	})

	return nil
}

func (s customerServer) StartBasket(
	ctx context.Context,
	request *pb_basket.StartBasketRequest,
) (*pb_basket.StartBasketResponse, error) {
	ctx = s.c.Scoped(ctx)

	var resp interface{}

	err := s.db.WithTx(ctx, func(ctx context.Context) (err error) {
		_, resp, err = s.startBasket.ServeGRPC(ctx, request)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp.(*pb_basket.StartBasketResponse), nil
}
