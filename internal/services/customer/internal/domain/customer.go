package domain

import (
	"github.com/htquangg/microservices-poc/internal/ddd"

	"github.com/stackus/errors"
)

const CustomerAggregate = "customers.CustomerAggregate"

type Customer struct {
	ddd.Aggregate
	name  string
	email string
	phone string
}

type CustomerOption func(*Customer) error

var (
	ErrNameCannotBeBlank       = errors.Wrap(errors.ErrBadRequest, "the customer username cannot be blank")
	ErrEmailAlreadyExists      = errors.Wrap(errors.ErrBadRequest, "the customer email is already existed")
	ErrCustomerIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the customer id cannot be blank")
	ErrSmsNumberCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "the SMS number cannot be blank")
	ErrCustomerNotAuthorized   = errors.Wrap(errors.ErrUnauthorized, "customer is not authorized")
)

func (c *Customer) Name() string {
	return c.name
}

func (c *Customer) Email() string {
	return c.email
}

func (c *Customer) Phone() string {
	return c.phone
}

func WithCustomerEmail(email string) CustomerOption {
	return func(c *Customer) error {
		c.email = email
		return nil
	}
}

func NewCustomer(id string) *Customer {
	return &Customer{
		Aggregate: ddd.NewAggregate(id, CustomerAggregate),
	}
}

func RegisterCustomer(id, name, phone string, options ...CustomerOption) (*Customer, error) {
	if id == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}
	if name == "" {
		return nil, ErrNameCannotBeBlank
	}

	customer := NewCustomer(id)
	customer.name = name
	customer.phone = phone

	for _, option := range options {
		if err := option(customer); err != nil {
			return nil, err
		}
	}

	customer.AddEvent(CustomerRegisteredEvent, &CustomerRegistered{
		Customer: customer,
	})

	return customer, nil
}
