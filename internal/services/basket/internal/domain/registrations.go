package domain

import (
	"github.com/htquangg/microservices-poc/internal/es"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/registry/serdes"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewJsonSerde(reg)

	// basket
	if err := serde.Register(BasketES{}, func(v interface{}) error {
		basket := v.(*BasketES)
		basket.Aggregate = es.NewAggregate("", BasketAggregate)
		basket.items = make(map[string]*Item)

		return nil
	}); err != nil {
		return err
	}

	// basket events
	if err := serde.Register(BasketStarted{}); err != nil {
		return err
	}
	if err := serde.Register(BasketCancelled{}); err != nil {
		return err
	}
	if err := serde.Register(BasketCheckedOut{}); err != nil {
		return err
	}
	if err := serde.Register(BasketItemAdded{}); err != nil {
		return err
	}
	if err := serde.Register(BasketItemRemoved{}); err != nil {
		return err
	}

	return nil
}
