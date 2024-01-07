package storepb

import (
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/registry/serdes"
)

const (
	StoreAggregateChannel = "mall.stores.events.Store"

	StoreCreatedEvent   = "storesapi.StoreCreated"
	StoreRebrandedEvent = "storesapi.StoreRebranded"

	ProductAggregateChannel = "mall.stores.events.Product"

	ProductAddedEvent = "storesapi.ProductAdded"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Store events
	if err := serde.Register(&StoreCreated{}); err != nil {
		return err
	}
	if err := serde.Register(&StoreRebranded{}); err != nil {
		return err
	}

	// Product events
	if err := serde.Register(&ProductAdded{}); err != nil {
		return err
	}

	return nil
}

func (*StoreCreated) Key() string {
	return StoreCreatedEvent
}

func (*StoreRebranded) Key() string {
	return StoreRebrandedEvent
}

func (*ProductAdded) Key() string {
	return ProductAddedEvent
}
