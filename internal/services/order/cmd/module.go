package cmd

import (
	"context"
	"time"

	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/es"
	"github.com/htquangg/microservices-poc/internal/kafka"
	mysql_internal "github.com/htquangg/microservices-poc/internal/mysql"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/services/basket/basketpb"
	"github.com/htquangg/microservices-poc/internal/services/order/constants"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/domain"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/grpc"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/handlers"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/system"
	"github.com/htquangg/microservices-poc/internal/services/order/orderpb"
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

			if err := domain.Registrations(reg); err != nil {
				return nil, err
			}
			if err := orderpb.Registrations(reg); err != nil {
				return nil, err
			}
			if err := basketpb.Registrations(reg); err != nil {
				return nil, err
			}

			return reg, nil
		},
	})

	builder.Add(di.Def{
		Name:  constants.DomainDispatcherKey,
		Scope: di.App,
		Build: func(_ di.Container) (interface{}, error) {
			return ddd.NewEventDispatcher[ddd.Event](), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.AggregateStoreKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return es.AggregateStoreWithMiddleware(
				mysql_internal.NewEventStore(
					svc.DB(),
					ctn.Get(constants.RegistryKey).(registry.Registry),
				),
				mysql_internal.NewSnapshotStore(
					svc.DB(),
					ctn.Get(constants.RegistryKey).(registry.Registry),
				),
			), nil
		},
	})
	kafkaProducer := kafka.NewProducer(svc.Config().Kafka.Brokers, svc.Logger())
	go func(ctx context.Context) error {
		<-ctx.Done()
		return kafkaProducer.Close()
	}(ctx)
	builder.Add(di.Def{
		Name:  constants.MessagePublisherKey,
		Scope: di.App,
		Build: func(_ di.Container) (interface{}, error) {
			outboxRepo := mysql_internal.NewOutboxStore(svc.DB())
			return am.NewMessagePublisher(kafkaProducer, tm.OutboxPublisher(outboxRepo)), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.InboxStoreKey,
		Scope: di.App,
		Build: func(_ di.Container) (interface{}, error) {
			return mysql_internal.NewInboxStore(svc.DB()), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.MessageSubscriberKey,
		Scope: di.App,
		Build: func(_ di.Container) (interface{}, error) {
			return am.NewMessageSubscriber(kafka.NewConsumer(&kafka.ConsumerConfig{
				Brokers:        svc.Config().Kafka.Brokers,
				Log:            svc.Logger(),
				Concurrency:    1,
				CommitInterval: time.Second,
			})), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.EventPublisherKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return am.NewEventPublisher(
				ctn.Get(constants.RegistryKey).(registry.Registry),
				ctn.Get(constants.MessagePublisherKey).(am.MessagePublisher),
			), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.OrderESRepoKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return es.NewAggregateRepository[*domain.OrderES](
				domain.OrderAggregate,
				ctn.Get(constants.RegistryKey).(registry.Registry),
				ctn.Get(constants.AggregateStoreKey).(es.AggregateStore),
			), nil
		},
	})

	// setup application
	builder.Add(di.Def{
		Name:  constants.ApplicationKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			domainDispatcher := ctn.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])
			return application.New(
					ctn.Get(constants.OrderESRepoKey).(domain.OrderESRepository),
					domainDispatcher,
					svc.Logger(),
				),
				nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.DomainEventHandlersKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return handlers.NewDomainEventHandlers(ctn.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.IntegrationEventHandlersKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return handlers.NewIntegrationEventHandlers(
					ctn.Get(constants.RegistryKey).(registry.Registry),
					ctn.Get(constants.ApplicationKey).(*application.Application),
					tm.InboxHandler(ctn.Get(constants.InboxStoreKey).(tm.InboxStore)),
				),
				nil
		},
	})
	outboxProcessor := tm.NewOutboxProcessor(kafkaProducer, mysql_internal.NewOutboxStore(svc.DB()))

	container := builder.Build()

	// setup driver adapters
	if err := grpc.RegisterServer(container, svc.DB(), svc.RPC()); err != nil {
		return err
	}
	handlers.RegisterDomainEventHandlers(container)
	handlers.RegisterIntegrationEventHandlers(container, svc.DB())
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
