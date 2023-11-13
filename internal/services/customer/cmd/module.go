package cmd

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/grpc"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/mysql"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/system"
)

func startUp(_ context.Context, svc system.Service) error {
	// setup Driven adapters
	customerRepo := mysql.NewCustomerRepository(svc.DB())
	domainDispatcher := ddd.NewEventDispatcher[ddd.AggregateEvent]()

	// setup application
	app := application.New(customerRepo, domainDispatcher, svc.Logger())

	// setup Driver adapters
	if err := grpc.RegisterServer(app, svc.DB(), svc.Sonyflake(), svc.RPC()); err != nil {
		return err
	}

	return nil
}
