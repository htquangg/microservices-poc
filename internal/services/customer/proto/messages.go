package proto

import (
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/registry/serdes"
)

const (
	CustomerAggregateChannel = "alex.customers.events.Customer"

	CustomerRegisteredEvent = "customerapi.CustomerRegistered"
	CustomerSmsChangedEvent = "customerapi.CustomerSmsChanged"
	CustomerEnabledEvent    = "customerapi.CustomerEnabled"
	CustomerDisabledEvent   = "customerapi.CustomerDisabled"

	CommandChannel = "alex.customers.commands"

	AuthorizeCustomerCommand = "customersapi.AuthorizeCustomer"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerde(reg)

	// Customer events
	if err := serde.Register(&CustomerRegistered{}); err != nil {
		return err
	}
	return nil
}

func (*CustomerRegistered) Key() string { return CustomerRegisteredEvent }
