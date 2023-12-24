package proto

import (
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/registry/serdes"
)

const (
	ProductAggregateChannel = "mall.stores.events.Product"

	ProductAddedEvent = "storesapi.ProductAdded"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	if err := serde.Register(&ProductAdded{}); err != nil {
		return err
	}

	return nil
}

func (*ProductAdded) Key() string { return ProductAddedEvent }
