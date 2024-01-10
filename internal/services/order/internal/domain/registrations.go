package domain

import (
	"github.com/htquangg/microservices-poc/internal/es"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/registry/serdes"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// basket
	if err := serde.Register(OrderES{}, func(v interface{}) error {
		basket := v.(*OrderES)
		basket.Aggregate = es.NewAggregate("", OrderAggregate)

		return nil
	}); err != nil {
		return err
	}

	// basket events
	if err := serde.Register(OrderCreated{}); err != nil {
		return err
	}
	if err := serde.Register(OrderRejected{}); err != nil {
		return err
	}
	if err := serde.Register(OrderApproved{}); err != nil {
		return err
	}
	if err := serde.Register(OrderCancelled{}); err != nil {
		return err
	}
	if err := serde.Register(OrderCompleted{}); err != nil {
		return err
	}

	return nil
}
