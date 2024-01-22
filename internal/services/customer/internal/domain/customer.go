package domain

import (
	"github.com/htquangg/microservices-poc/internal/ddd"

	"github.com/stackus/errors"
)

const CustomerAggregate = "customers.CustomerAggregate"

type Customer struct {
	ddd.Aggregate
	name  string
	phone string
	email string
}

type CustomerOption func(*Customer) error

var (
	ErrCustomerIDCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "the customer id cannot be blank")
	ErrNameCannotBeBlank        = errors.Wrap(errors.ErrBadRequest, "the customer username cannot be blank")
	ErrPhoneNumberCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the phone number cannot be blank")
	ErrEmailCannotBeBlank       = errors.Wrap(errors.ErrBadRequest, "the email cannot be blank")
	ErrEmailAlreadyExists       = errors.Wrap(errors.ErrBadRequest, "the customer email is already existed")
	ErrCustomerNotAuthorized    = errors.Wrap(errors.ErrUnauthorized, "customer is not authorized")
)

func (c *Customer) Name() string {
	return c.name
}

func (c *Customer) Phone() string {
	return c.phone
}

func (c *Customer) Email() string {
	return c.email
}

func NewCustomer(id string, name, phone, email string, options ...CustomerOption) (*Customer, error) {
	if id == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}

	if name == "" {
		return nil, ErrNameCannotBeBlank
	}

	if phone == "" {
		return nil, ErrPhoneNumberCannotBeBlank
	}

	if email == "" {
		return nil, ErrEmailCannotBeBlank
	}

	customer := &Customer{
		Aggregate: ddd.NewAggregate(id, CustomerAggregate),
	}

	for _, option := range options {
		if err := option(customer); err != nil {
			return nil, err
		}
	}

	return customer, nil
}

func RegisterCustomer(id, name, phone, email string, options ...CustomerOption) (*Customer, error) {
	customer, err := NewCustomer(id, name, phone, email, options...)
	if err != nil {
		return nil, err
	}

	customer.AddEvent(CustomerRegisteredEvent, &CustomerRegistered{
		Customer: customer,
	})

	return customer, nil
}

func (c *Customer) Authorize() error {
	c.AddEvent(CustomerAuthorizedEvent, &CustomerAuthorized{
		Customer: c,
	})

	return nil
}
