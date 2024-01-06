package proto

import (
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/registry/serdes"
)

const (
	BasketAggregateChannel = "mall.stores.events.Basket"

	BasketStartedEvent    = "basketapi.BasketStarted"
	BasketCancelledEvent  = "basketapi.BasketCancelled"
	BasketCheckedOutEvent = "basketapi.BasketCheckedOut"
)

func Registrations(reg registry.Registry) error {
	return RegistrationsWithSerde(serdes.NewProtoSerde(reg))
}

func RegistrationsWithSerde(serde registry.Serde) error {
	// Basket events
	if err := serde.Register(&BasketStarted{}); err != nil {
		return err
	}
	if err := serde.Register(&BasketCancelled{}); err != nil {
		return err
	}
	if err := serde.Register(&BasketCheckedOut{}); err != nil {
		return err
	}

	return nil
}

func (*BasketStarted) Key() string {
	return BasketStartedEvent
}

func (*BasketCancelled) Key() string {
	return BasketCancelledEvent
}

func (*BasketCheckedOut) Key() string {
	return BasketCheckedOutEvent
}
