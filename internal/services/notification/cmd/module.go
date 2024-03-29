package cmd

import (
	"context"
	"time"

	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/kafka"
	mysql_internal "github.com/htquangg/microservices-poc/internal/mysql"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/services/customer/customerpb"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/grpc"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/handlers"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/mysql"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/system"
	"github.com/htquangg/microservices-poc/internal/tm"
)

func startUp(ctx context.Context, svc system.Service) error {
	// setup driven adapters
	reg := registry.New()
	if err := customerpb.Registrations(reg); err != nil {
		return err
	}
	inboxStore := mysql_internal.NewInboxStore(svc.DB())
	messageSubscriber := am.NewMessageSubscriber(kafka.NewConsumer(&kafka.ConsumerConfig{
		Brokers:        svc.Config().Kafka.Brokers,
		Log:            svc.Logger(),
		Concurrency:    1,
		CommitInterval: time.Second,
	}))
	customerRepo := mysql.NewCustomerRepository(svc.DB())

	// setup application
	app := application.New()
	integrationEventHandlers := handlers.NewIntegrationEventHandlers(reg,
		app,
		customerRepo,
		tm.InboxHandler(inboxStore),
	)

	// setup Driver adapters
	if err := grpc.RegisterServer(ctx, app, svc.RPC()); err != nil {
		return err
	}
	if err := handlers.RegisterIntegrationEventHandlers(messageSubscriber, integrationEventHandlers); err != nil {
		return err
	}

	return nil
}
