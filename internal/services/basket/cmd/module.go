package cmd

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/services/basket/constants"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/grpc"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/system"

	"github.com/htquangg/di/v2"
)

func startUp(_ context.Context, svc system.Service) error {
	builder, err := di.NewBuilder()
	if err != nil {
		return err
	}

	// setup driven adapters
	builder.Add(di.Def{
		Name:  constants.RegistryKey,
		Scope: di.App,
		Build: func(_ di.Container) (interface{}, error) {
			reg := registry.New()
			return reg, nil
		},
	})

	// setup application
	container := builder.Build()

	// setup driver adapters
	if err := grpc.RegisterServer(container, svc.DB(), svc.RPC()); err != nil {
		return err
	}

	return nil
}
