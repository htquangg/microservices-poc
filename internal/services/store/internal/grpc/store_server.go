package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/store/storepb"
	"github.com/htquangg/microservices-poc/pkg/database"

	grpc_transport "github.com/go-kit/kit/transport/grpc"
	"github.com/htquangg/di/v2"
	"google.golang.org/grpc"
)

var _ storepb.StoreServiceServer = (*storeServer)(nil)

type storeServer struct {
	storepb.UnimplementedStoreServiceServer

	ctn di.Container
	db  database.DB

	createStore  grpc_transport.Handler
	rebrandStore grpc_transport.Handler
	addProduct   grpc_transport.Handler
}

func registerStoreServer(
	ctn di.Container,
	db database.DB,
	registrar grpc.ServiceRegistrar,
) error {
	endpoints := makeStoreEndpoints(ctn)

	storepb.RegisterStoreServiceServer(registrar, storeServer{
		ctn: ctn,
		db:  db,
		createStore: grpc_transport.NewServer(
			endpoints.createStoreEndpoint,
			decodeCreateStoreRequest,
			encodeCreateStoreResponse,
		),
		rebrandStore: grpc_transport.NewServer(
			endpoints.rebrandStoreEndpoint,
			decodeRebrandStoreRequest,
			encodeRebrandStoreResponse,
		),
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
	request *storepb.AddProductRequest,
) (*storepb.AddProductResponse, error) {
	var resp interface{}

	err := s.db.WithTx(ctx, func(ctx context.Context) (err error) {
		_, resp, err = s.addProduct.ServeGRPC(ctx, request)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp.(*storepb.AddProductResponse), nil
}

func (s storeServer) CreateStore(
	ctx context.Context,
	request *storepb.CreateStoreRequest,
) (*storepb.CreateStoreResponse, error) {
	var resp interface{}

	err := s.db.WithTx(ctx, func(ctx context.Context) (err error) {
		_, resp, err = s.createStore.ServeGRPC(ctx, request)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp.(*storepb.CreateStoreResponse), nil
}

func (s storeServer) RebrandStore(
	ctx context.Context,
	request *storepb.RebrandStoreRequest,
) (*storepb.RebrandStoreResponse, error) {
	var resp interface{}

	err := s.db.WithTx(ctx, func(ctx context.Context) (err error) {
		_, resp, err = s.rebrandStore.ServeGRPC(ctx, request)
		return err
	})
	if err != nil {
		return nil, err
	}

	return resp.(*storepb.RebrandStoreResponse), nil
}
