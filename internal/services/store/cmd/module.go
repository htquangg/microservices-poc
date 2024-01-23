package cmd

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/es"
	"github.com/htquangg/microservices-poc/internal/kafka"
	mysql_internal "github.com/htquangg/microservices-poc/internal/mysql"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/registry/serdes"
	"github.com/htquangg/microservices-poc/internal/services/store/constants"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/domain"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/grpc"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/handlers"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/mysql"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/system"
	"github.com/htquangg/microservices-poc/internal/services/store/storepb"
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
			if err := registrations(reg); err != nil {
				return nil, err
			}
			if err := storepb.Registrations(reg); err != nil {
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
	builder.Add(di.Def{
		Name:  constants.StoreRepoKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return mysql.NewStoreRepository(svc.DB()), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.StoreESRepoKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return es.NewAggregateRepository[*domain.StoreES](
				domain.StoreAggregate,
				ctn.Get(constants.RegistryKey).(registry.Registry),
				ctn.Get(constants.AggregateStoreKey).(es.AggregateStore),
			), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.ProductRepoKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return mysql.NewProductRepository(svc.DB()), nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.ProductESRepoKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return es.NewAggregateRepository[*domain.ProductES](
				domain.ProductAggregate,
				ctn.Get(constants.RegistryKey).(registry.Registry),
				ctn.Get(constants.AggregateStoreKey).(es.AggregateStore),
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
		Name:  constants.EventPublisherKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return am.NewEventPublisher(
				ctn.Get(constants.RegistryKey).(registry.Registry),
				ctn.Get(constants.MessagePublisherKey).(am.MessagePublisher),
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
					ctn.Get(constants.StoreESRepoKey).(domain.StoreESRepository),
					ctn.Get(constants.ProductESRepoKey).(domain.ProductESRepository),
					domainDispatcher,
					svc.Logger(),
				),
				nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.StoreHandlersKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return handlers.NewStoreHandlers(
					ctn.Get(constants.StoreRepoKey).(domain.StoreRepository),
				),
				nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.ProductHandlersKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return handlers.NewProductHandlers(
					ctn.Get(constants.ProductRepoKey).(domain.ProductRepository),
				),
				nil
		},
	})
	builder.Add(di.Def{
		Name:  constants.DomainEventHandlersKey,
		Scope: di.App,
		Build: func(ctn di.Container) (interface{}, error) {
			return handlers.NewDomainEventHandlers(
					ctn.Get(constants.EventPublisherKey).(am.EventPublisher),
				),
				nil
		},
	})
	outboxProcessor := tm.NewOutboxProcessor(kafkaProducer, mysql_internal.NewOutboxStore(svc.DB()))

	ctn := builder.Build()

	// setup driver adapters
	if err := grpc.RegisterServer(ctn, svc.DB(), svc.RPC()); err != nil {
		return err
	}
	handlers.RegisterDomainEventHandlers(ctn)
	handlers.RegisterStoreHandlers(ctn)
	handlers.RegisterProductHandlers(ctn)
	startOutboxProcessor(ctx, outboxProcessor, svc.Logger())

	return nil
}

func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// store
	if err := serde.Register(domain.StoreES{}, func(v any) error {
		store := v.(*domain.StoreES)
		store.Aggregate = es.NewAggregate("", domain.StoreAggregate)
		return nil
	}); err != nil {
		return err
	}

	// store events
	if err = serde.Register(domain.StoreCreated{}); err != nil {
		return
	}
	if err = serde.Register(domain.StoreRebranded{}); err != nil {
		return
	}

	// store snapshots
	if err = serde.RegisterKey(domain.StoreV1{}.SnapshotName(), domain.StoreV1{}); err != nil {
		return
	}

	// product
	if err = serde.Register(domain.ProductES{}, func(v any) error {
		store := v.(*domain.ProductES)
		store.Aggregate = es.NewAggregate("", domain.ProductAggregate)
		return nil
	}); err != nil {
		return
	}

	// product events
	if err = serde.Register(domain.ProductAdded{}); err != nil {
		return
	}

	// product snapshots
	if err = serde.RegisterKey(domain.ProductV1{}.SnapshotName(), domain.ProductV1{}); err != nil {
		return
	}

	return
}

func startOutboxProcessor(ctx context.Context, outboxProcessor tm.OutboxProcessor, log logger.Logger) {
	go func() {
		err := outboxProcessor.Start(ctx)
		if err != nil {
			log.Err("store outbox processor encountered an error", err)
		}
	}()
}
