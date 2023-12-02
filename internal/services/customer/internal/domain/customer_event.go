package domain

const (
	CustomerRegisteredEvent = "customers.CustomerRegistered"
)

type CustomerRegistered struct {
	Customer *Customer
}

func (CustomerRegistered) Key() string {
	return CustomerRegisteredEvent
}
