package grpc

import (
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application"
	"github.com/htquangg/microservices-poc/pkg/database"
	"github.com/htquangg/microservices-poc/pkg/uid"

	"google.golang.org/grpc"
)

type server struct {
	app *application.Application
	db  *database.DB
	sf  *uid.Sonyflake
}

func RegisterServer(
	app *application.Application,
	db *database.DB,
	sf *uid.Sonyflake,
	registrar grpc.ServiceRegistrar,
) error {
	if err := registerCustomerServer(app, db, sf, registrar); err != nil {
		return err
	}

	return nil
}
