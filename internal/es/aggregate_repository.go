package es

import (
	"context"
	"fmt"

	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/registry"
)

type (
	AggreagateRepository[T EventSourcedAggregate] interface {
		Load(ctx context.Context, aggregateID string) (agg T, err error)
		Save(ctx context.Context, aggregate T) error
	}

	aggreagateRepository[T EventSourcedAggregate] struct {
		aggregateName string
		registry      registry.Registry
		store         AggregateStore
	}
)

var _ AggreagateRepository[EventSourcedAggregate] = (*aggreagateRepository[EventSourcedAggregate])(nil)

func NewAggregateRepository[T EventSourcedAggregate](
	aggreagateName string,
	registry registry.Registry,
	store AggregateStore,
) *aggreagateRepository[T] {
	return &aggreagateRepository[T]{
		aggregateName: aggreagateName,
		registry:      registry,
		store:         store,
	}
}

func (r *aggreagateRepository[T]) Load(ctx context.Context, aggregateID string) (agg T, err error) {
	var v interface{}
	v, err = r.registry.Build(
		r.aggregateName,
		ddd.SetID(aggregateID),
		ddd.SetName(r.aggregateName),
	)
	if err != nil {
		return agg, err
	}

	var ok bool
	if agg, ok = v.(T); !ok {
		return agg, fmt.Errorf("%T is not the expected type %T", v, agg)
	}

	if err = r.store.Load(ctx, agg); err != nil {
		return agg, err
	}

	return agg, nil
}

func (r *aggreagateRepository[T]) Save(ctx context.Context, aggregate T) error {
	if aggregate.Version() == aggregate.PendingVersion() {
		return nil
	}

	for _, event := range aggregate.Events() {
		if err := aggregate.ApplyEvent(event); err != nil {
			return err
		}
	}

	if err := r.store.Save(ctx, aggregate); err != nil {
		return err
	}

	aggregate.CommitEvents()

	return nil
}
