package domain

const (
	CustomerRegisteredEvent   = "customers.CustomerRegistered"
	CustomerPhoneChangedEvent = "customers.CustomerPhoneChanged"
)

type CustomerRegistered struct {
	Customer *Customer
}

func (CustomerRegistered) Key() string {
	return CustomerRegisteredEvent
}

type CustomerPhoneChanged struct {
	Customer *Customer
}

func (CustomerPhoneChanged) Key() string {
	return CustomerPhoneChangedEvent
}
