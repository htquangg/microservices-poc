package customerpb

import (
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/registry/serdes"
)

const (
	CustomerAggregateChannel = "mall.customers.events.Customer"

	CustomerRegisteredEvent = "customerapi.CustomerRegistered"
	CustomerSmsChangedEvent = "customerapi.CustomerSmsChanged"
	CustomerEnabledEvent    = "customerapi.CustomerEnabled"
	CustomerDisabledEvent   = "customerapi.CustomerDisabled"

	CommandChannel = "mall.customers.commands"

	AuthorizeCustomerCommand = "customerapi.AuthorizeCustomer"
)

func Registrations(reg registry.Registry) error {
	serde := serdes.NewProtoSerde(reg)

	// events
	if err := serde.Register(&CustomerRegistered{}); err != nil {
		return err
	}

	// commands
	if err := serde.Register(&AuthorizeCustomer{}); err != nil {
		return err
	}

	return nil
}

func (*CustomerRegistered) Key() string {
	return CustomerRegisteredEvent
}

func (*AuthorizeCustomer) Key() string {
	return AuthorizeCustomerCommand
}
