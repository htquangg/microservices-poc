package orderpb

import (
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/registry/serdes"
)

const (
	OrderAggregateChannel = "mallbots.ordering.events.Order"

	OrderCreatedEvent   = "ordersapi.OrderCreated"
	OrderRejectedEvent  = "ordersapi.OrderRejected"
	OrderApprovedEvent  = "ordersapi.OrderApproved"
	OrderCanceledEvent  = "ordersapi.OrderCanceled"
	OrderCompletedEvent = "ordersapi.OrderCompleted"
)

func Registrations(reg registry.Registry) (err error) {
	serde := serdes.NewProtoSerde(reg)

	// Order events
	if err = serde.Register(&OrderCreated{}); err != nil {
		return err
	}
	if err = serde.Register(&OrderRejected{}); err != nil {
		return err
	}
	if err = serde.Register(&OrderApproved{}); err != nil {
		return err
	}
	if err = serde.Register(&OrderCanceled{}); err != nil {
		return err
	}
	if err = serde.Register(&OrderCompleted{}); err != nil {
		return err
	}

	return nil
}

func (*OrderCreated) Key() string {
	return OrderCreatedEvent
}

func (*OrderRejected) Key() string {
	return OrderRejectedEvent
}

func (*OrderApproved) Key() string {
	return OrderApprovedEvent
}

func (*OrderCanceled) Key() string {
	return OrderCanceledEvent
}

func (*OrderCompleted) Key() string {
	return OrderCompletedEvent
}