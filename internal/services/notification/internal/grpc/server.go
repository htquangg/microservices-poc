package grpc

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/notification/internal/application"
	notificationpb "github.com/htquangg/microservices-poc/internal/services/notification/proto"

	"google.golang.org/grpc"
)

type server struct {
	app *application.Application
	notificationpb.UnimplementedNotificationServiceServer
}

func RegisterServer(
	_ context.Context,
	app *application.Application,
	registrar grpc.ServiceRegistrar,
) error {
	notificationpb.RegisterNotificationServiceServer(registrar, server{app: app})
	return nil
}
