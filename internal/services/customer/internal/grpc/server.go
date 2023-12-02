package grpc

import (
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application"
	"github.com/htquangg/microservices-poc/pkg/database"
	"github.com/htquangg/microservices-poc/pkg/uid"

	"github.com/htquangg/di/v2"
	"google.golang.org/grpc"
)

type server struct {
	c   di.Container
	app *application.Application
	db  *database.DB
	sf  *uid.Sonyflake
}

func RegisterServer(
	c di.Container,
	db *database.DB,
	sf *uid.Sonyflake,
	registrar grpc.ServiceRegistrar,
) error {
	if err := registerCustomerServer(c, db, sf, registrar); err != nil {
		return err
	}

	return nil
}
