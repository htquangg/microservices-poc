package cmd

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/kafka"
	mysql_internal "github.com/htquangg/microservices-poc/internal/mysql"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/services/customer/constants"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/grpc"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/handlers"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/mysql"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/system"
	pb_customer "github.com/htquangg/microservices-poc/internal/services/customer/proto"
	"github.com/htquangg/microservices-poc/internal/tm"
	"github.com/htquangg/microservices-poc/pkg/logger"

	"github.com/htquangg/di/v2"
)

func startUp(ctx context.Context, svc system.Service) error {
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
			if err := pb_customer.Registrations(reg); err != nil {
				return nil, err
			}
			return reg, nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.CustomerRepoKey,
		Scope: di.Request,
		Build: func(c di.Container) (interface{}, error) {
			return mysql.NewCustomerRepository(svc.DB()), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.DomainDispatcherKey,
		Scope: di.App,
		Build: func(_ di.Container) (interface{}, error) {
			return ddd.NewEventDispatcher[ddd.AggregateEvent](), nil
		},
	})
	kafkaProducer := kafka.NewProducer(svc.Config().Kafka.Brokers, svc.Logger())
	go func(ctx context.Context) error {
		<-ctx.Done()
		return kafkaProducer.Close()
	}(ctx)
	builder.Add(di.Def{
		Name:  constants.MessagePublisherKey,
		Scope: di.Request,
		Build: func(_ di.Container) (interface{}, error) {
			outboxRepo := mysql_internal.NewOutboxStore(svc.DB())
			return am.NewMessagePublisher(kafkaProducer, tm.OutboxPublisher(outboxRepo)), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.EventPublisherKey,
		Scope: di.Request,
		Build: func(c di.Container) (interface{}, error) {
			return am.NewEventPublisher(
				c.Get(constants.RegistryKey).(registry.Registry),
				c.Get(constants.MessagePublisherKey).(am.MessagePublisher),
			), nil
		},
	})

	// setup application
	builder.Add(di.Def{
		Name:  constants.ApplicationKey,
		Scope: di.Request,
		Build: func(c di.Container) (interface{}, error) {
			domainDispatcher := c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.AggregateEvent])
			return application.New(
				c.Get(constants.CustomerRepoKey).(mysql.CustomerRepository),
				domainDispatcher,
				svc.Logger(),
			), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.DomainEventHandlersKey,
		Scope: di.Request,
		Build: func(c di.Container) (interface{}, error) {
			return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
		},
	})
	outboxProcessor := tm.NewOutboxProcessor(kafkaProducer, mysql_internal.NewOutboxStore(svc.DB()))

	container := builder.Build()

	// setup driver adapters
	if err := grpc.RegisterServer(container, svc.DB(), svc.RPC()); err != nil {
		return err
	}
	handlers.RegisterDomainEventHandlers(container)
	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, log logger.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			log.Err("customer outbox processor encountered an error", err)
		}
	}()
}
