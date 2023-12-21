package cmd

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/es"
	mysql_internal "github.com/htquangg/microservices-poc/internal/mysql"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/registry/serdes"
	"github.com/htquangg/microservices-poc/internal/services/store/constants"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/domain"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/grpc"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/system"
	pb_store "github.com/htquangg/microservices-poc/internal/services/store/proto"

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
			if err := registrations(reg); err != nil {
				return nil, err
			}
			if err := pb_store.Registrations(reg); err != nil {
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
	builder.Add(di.Def{
		Name:  constants.ProductRepoKey,
		Scope: di.Request,
		Build: func(c di.Container) (interface{}, error) {
			return es.NewAggregateRepository[*domain.Product](
				domain.ProductAggregate,
				c.Get(constants.RegistryKey).(registry.Registry),
				c.Get(constants.AggregateStoreKey).(es.AggregateStore),
			), nil
		},
	})

	// setup application
	builder.Add(di.Def{
		Name:  constants.ApplicationKey,
		Scope: di.Request,
		Build: func(c di.Container) (interface{}, error) {
			domainDispatcher := c.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])
			return application.New(
				c.Get(constants.ProductRepoKey).(domain.ProductRepository),
				domainDispatcher,
				svc.Logger(),
			), nil
		},
	})

	container := builder.Build()

	// setup driver adapters
	if err := grpc.RegisterServer(container, svc.DB(), svc.RPC()); err != nil {
		return err
	}

	return nil
}

func registrations(reg registry.Registry) (err error) {
	serde := serdes.NewJsonSerde(reg)

	// Product
	if err = serde.Register(domain.Product{}, func(v any) error {
		store := v.(*domain.Product)
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
