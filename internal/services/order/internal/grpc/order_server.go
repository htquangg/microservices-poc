package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/order/orderpb"
	"github.com/htquangg/microservices-poc/pkg/database"

	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"github.com/htquangg/di/v2"
	"google.golang.org/grpc"
)

var _ orderpb.OrderServiceServer = (*orderServer)(nil)

type orderServer struct {
	orderpb.UnimplementedOrderServiceServer

	ctn di.Container
	db  database.DB

	createOrder grpc_transport.Handler
	cancelOrder grpc_transport.Handler
}

func registerOrderServer(
	ctn di.Container,
	db database.DB,
	registrar grpc.ServiceRegistrar,
) error {
	endpoints := makeOrderEndpoints(ctn)

	orderpb.RegisterOrderServiceServer(registrar, orderServer{
		ctn: ctn,
		db:  db,
		createOrder: grpc_transport.NewServer(
			endpoints.createOrderEndpoint,
			decodeCreateOrderRequest,
			endcodeCreateOrderResponse,
		),
		cancelOrder: grpc_transport.NewServer(
			endpoints.cancelOrderEndpoint,
			decodeCancelOrderRequest,
			encodeCancelOrderResponse,
		),
	})

	return nil
}

func (s orderServer) CreateOrder(
	ctx context.Context,
	request *orderpb.CreateOrderRequest,
) (*orderpb.CreateOrderResponse, error) {
	var resp interface{}

	err := s.db.WithTx(ctx, func(ctx context.Context) (err error) {
		_, resp, err = s.createOrder.ServeGRPC(ctx, request)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp.(*orderpb.CreateOrderResponse), nil
}

func (s orderServer) CancelOrder(
	ctx context.Context,
	request *orderpb.CancelOrderRequest,
) (*orderpb.CancelOrderResponse, error) {
	var resp interface{}

	err := s.db.WithTx(ctx, func(ctx context.Context) (err error) {
		_, resp, err = s.cancelOrder.ServeGRPC(ctx, request)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp.(*orderpb.CancelOrderResponse), nil
}
