package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application/command"
	customerpb "github.com/htquangg/microservices-poc/internal/services/customer/proto"
	"github.com/htquangg/microservices-poc/pkg/database"
	"github.com/htquangg/microservices-poc/pkg/uid"

	"google.golang.org/grpc"
)

var _ customerpb.CustomerServiceServer = (*server)(nil)

type server struct {
	app *application.Application
	db  *database.DB
	sf  *uid.Sonyflake
	customerpb.UnimplementedCustomerServiceServer
}

func RegisterServer(
	app *application.Application,
	db *database.DB,
	sf *uid.Sonyflake,
	registrar grpc.ServiceRegistrar,
) error {
	customerpb.RegisterCustomerServiceServer(registrar, &server{
		app: app,
		db:  db,
		sf:  sf,
	})
	return nil
}

func (s *server) RegisterCustomer(
	ctx context.Context,
	request *customerpb.RegisterCustomerRequest,
) (*customerpb.RegisterCustomerResponse, error) {
	id := s.sf.ID()
	errTx := s.db.WithTx(ctx, func(ctx context.Context) error {
		err := s.app.Commands.RegisterCustomerHandler.Handle(ctx, command.RegisterCustomer{
			ID:    id,
			Name:  request.GetName(),
			Phone: request.GetPhone(),
		})
		return err
	})
	return &customerpb.RegisterCustomerResponse{Id: id}, errTx
}
