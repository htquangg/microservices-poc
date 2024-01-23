package grpc

import (
	"github.com/htquangg/microservices-poc/pkg/database"

	"github.com/htquangg/di/v2"
	"google.golang.org/grpc"
)

func RegisterServer(
	ctn di.Container,
	db database.DB,
	registrar grpc.ServiceRegistrar,
) error {
	return nil
}
