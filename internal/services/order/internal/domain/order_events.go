package domain

const (
	OrderCreatedEvent   = "orders.OrderCreated"
	OrderRejectedEvent  = "orders.OrderRejected"
	OrderApprovedEvent  = "orders.OrderApproved"
	OrderCancelledEvent = "orders.OrderCancelled"
	OrderCompletedEvent = "orders.OrderCompleted"
)

type OrderCreated struct {
	CustomerID string
	PaymentID  string
	ShoppingID string
	Items      []*Item
}

// Key implements registry.Registrable
func (OrderCreated) Key() string {
	return OrderCreatedEvent
}

type OrderRejected struct{}

// Key implements registry.Registrable
func (OrderRejected) Key() string {
	return OrderRejectedEvent
}

type OrderApproved struct {
	ShoppingID string
}

// Key implements registry.Registrable
func (OrderApproved) Key() string {
	return OrderApprovedEvent
}

type OrderCancelled struct {
	CustomerID string
	PaymentID  string
}

// Key implements registry.Registrable
func (OrderCancelled) Key() string {
	return OrderCancelledEvent
}

type OrderCompleted struct {
	CustomerID string
	InvoiceID  string
}

// Key implements registry.Registrable
func (OrderCompleted) Key() string {
	return OrderCompletedEvent
}
