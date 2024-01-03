package domain

const (
	BasketStartedEvent     = "baskets.BasketStarted"
	BasketItemAddedEvent   = "baskets.BasketItemAdded"
	BasketItemRemovedEvent = "baskets.BasketItemRemoved"
	BasketCancelledEvent   = "baskets.BasketCancelled"
	BasketCheckedOutEvent  = "baskets.BasketCheckedOut"
)

type BasketStarted struct {
	CustomerID string
}

// Key implements registry.Registrable
func (BasketStarted) Key() string {
	return BasketStartedEvent
}

type BasketItemAdded struct {
	Item Item
}

// Key implements registry.Registrable
func (BasketItemAdded) Key() string {
	return BasketItemAddedEvent
}

type BasketItemRemoved struct {
	ProductID string
	Quantity  int
}

// Key implements registry.Registrable
func (BasketItemRemoved) Key() string {
	return BasketItemRemovedEvent
}

type BasketCancelled struct{}

// Key implements registry.Registrable
func (BasketCancelled) Key() string {
	return BasketCancelledEvent
}

type BasketCheckedOut struct {
	PaymentID string
}

// Key implements registry.Registrable
func (BasketCheckedOut) Key() string {
	return BasketCheckedOutEvent
}
