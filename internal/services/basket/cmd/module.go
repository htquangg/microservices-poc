package cmd

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/es"
	"github.com/htquangg/microservices-poc/internal/kafka"
	mysql_internal "github.com/htquangg/microservices-poc/internal/mysql"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/services/basket/constants"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/domain"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/grpc"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/handlers"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/mysql"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/system"
	pb_basket "github.com/htquangg/microservices-poc/internal/services/basket/proto"
	"github.com/htquangg/microservices-poc/internal/tm"

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
			if err := pb_basket.Registrations(reg); err != nil {
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
		Scope: di.Request,
		Build: func(c di.Container) (interface{}, error) {
			return es.AggregateStoreWithMiddleware(
				mysql_internal.NewEventStore(
					svc.DB(),
					c.Get(constants.RegistryKey).(registry.Registry),
				),
				mysql_internal.NewSnapshotStore(
					svc.DB(),
					c.Get(constants.RegistryKey).(registry.Registry),
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
	builder.Add(di.Def{
		Name:  constants.BasketESRepoKey,
		Scope: di.Request,
		Build: func(c di.Container) (interface{}, error) {
			return es.NewAggregateRepository[*domain.BasketES](
				domain.BasketAggregate,
				c.Get(constants.RegistryKey).(registry.Registry),
				c.Get(constants.AggregateStoreKey).(es.AggregateStore),
			), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.StoreRepoKey,
		Scope: di.Request,
		Build: func(c di.Container) (interface{}, error) {
			return mysql.NewStoreRepository(svc.DB(), nil), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.ProductRepoKey,
		Scope: di.Request,
		Build: func(c di.Container) (interface{}, error) {
			return mysql.NewProductRepository(svc.DB(), nil), nil
		},
	})

	// setup application
	builder.Add(di.Def{
		Name:  constants.ApplicationKey,
		Scope: di.Request,
		Build: func(c di.Container) (interface{}, error) {
			domainDispatcher := c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])
			return application.New(
					c.Get(constants.BasketESRepoKey).(domain.BasketESRepository),
					c.Get(constants.StoreRepoKey).(domain.StoreRepository),
					c.Get(constants.ProductRepoKey).(domain.ProductRepository),
					domainDispatcher,
					svc.Logger(),
				),
				nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.DomainEventHandlersKey,
		Scope: di.Request,
		Build: func(c di.Container) (interface{}, error) {
			return handlers.NewDomainEventHandlers(c.Get(constants.EventPublisherKey).(am.EventPublisher)), nil
		},
	})

	container := builder.Build()

	// setup driver adapters
	if err := grpc.RegisterServer(container, svc.DB(), svc.RPC()); err != nil {
		return err
	}
	handlers.RegisterDomainEventHandlers(container)

	return nil
}
