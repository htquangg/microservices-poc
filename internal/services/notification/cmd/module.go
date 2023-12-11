package cmd

import (
	"context"
	"time"

	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/kafka"
	"github.com/htquangg/microservices-poc/internal/registry"
	customerpb "github.com/htquangg/microservices-poc/internal/services/customer/proto"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/grpc"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/handlers"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/mysql"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/system"
)

func startUp(ctx context.Context, svc system.Service) error {
	// setup driven adapters
	reg := registry.New()
	if err := customerpb.Registrations(reg); err != nil {
		return err
	}
	messageSubscriber := am.NewMessageSubscriber(kafka.NewConsumer(&kafka.ConsumerConfig{
		Brokers:        svc.Config().Kafka.Brokers,
		Log:            svc.Logger(),
		Concurrency:    1,
		CommitInterval: time.Second,
	}))
	customerRepo := mysql.NewCustomerRepository(svc.DB())

	// setup application
	app := application.New()
	intergrationEventHanders := handlers.NewIntegrationEventHandlers(reg,
		app,
		customerRepo,
	)

	// setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, svc.RPC()); err != nil {
		return err
	}
	if err := handlers.RegisterIntergrationEventHandlers(messageSubscriber, intergrationEventHanders); err != nil {
		return err
	}

	return nil
}
