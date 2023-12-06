package cmd

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/registry"
	customerpb "github.com/htquangg/microservices-poc/internal/services/customer/proto"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/grpc"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/system"
)

func startUp(ctx context.Context, svc system.Service) error {
	// setup driven adapters
	reg := registry.New()
	if err := customerpb.Registrations(reg); err != nil {
		return err
	}

	// setup application
	app := application.New()

	// setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, svc.RPC()); err != nil {
		return err
	}

	return nil
}
