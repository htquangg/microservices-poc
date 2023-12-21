package domain

const (
	CustomerRegisteredEvent = "customers.CustomerRegistered"
)

type CustomerRegistered struct {
	Customer *Customer
}

// Key implements registry.Registerable
func (CustomerRegistered) Key() string {
	return CustomerRegisteredEvent
}
