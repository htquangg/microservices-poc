package domain

const (
	CustomerRegisteredEvent = "customers.CustomerRegistered"
	CustomerAuthorizedEvent = "customers.CustomerAuthorized"
)

type CustomerRegistered struct {
	Customer *Customer
}

// Key implements registry.Registerable
func (CustomerRegistered) Key() string {
	return CustomerRegisteredEvent
}

type CustomerAuthorized struct {
	Customer *Customer
}

// Key implements registry.Registerable
func (CustomerAuthorized) Key() string {
	return CustomerAuthorizedEvent
}
