package grpc

import (
	"context"

	pb_store "github.com/htquangg/microservices-poc/internal/services/store/proto"
	"github.com/htquangg/microservices-poc/pkg/database"

	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"github.com/htquangg/di/v2"
	"google.golang.org/grpc"
)

var _ pb_store.StoreServiceServer = (*storeServer)(nil)

type storeServer struct {
	pb_store.UnimplementedStoreServiceServer

	c  di.Container
	db database.DB

	addProduct grpc_transport.Handler
}

func registerStoreServer(
	c di.Container,
	db database.DB,
	registrar grpc.ServiceRegistrar,
) error {
	endpoints := makeStoreEndpoints(c)

	pb_store.RegisterStoreServiceServer(registrar, storeServer{
		c:  c,
		db: db,
		addProduct: grpc_transport.NewServer(
			endpoints.registerCustomerEndpoint,
			decodeAddProductRequest,
			encodeAddProductResponse,
		),
	})

	return nil
}

func (s storeServer) AddProduct(
	ctx context.Context,
	request *pb_store.AddProductRequest,
) (*pb_store.AddProductResponse, error) {
	ctx = s.c.Scoped(ctx)

	var resp interface{}

	err := s.db.WithTx(ctx, func(ctx context.Context) (err error) {
		_, resp, err = s.addProduct.ServeGRPC(ctx, request)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp.(*pb_store.AddProductResponse), nil
}
