package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/notification/internal/application"
	pb_notification "github.com/htquangg/microservices-poc/internal/services/notification/proto"

	"google.golang.org/grpc"
)

type server struct {
	app *application.Application
	pb_notification.UnimplementedNotificationServiceServer
}

func RegisterServer(
	_ context.Context,
	app *application.Application,
	registrar grpc.ServiceRegistrar,
) error {
	pb_notification.RegisterNotificationServiceServer(registrar, server{app: app})
	return nil
}
