package domain

import (
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/es"

	"github.com/stackus/errors"
)

const BasketAggregate = "baskets.Basket"

var (
	ErrBasketHasNoItems               = errors.Wrap(errors.ErrBadRequest, "the basket has no items")
	ErrBasketCannotBeModified         = errors.Wrap(errors.ErrBadRequest, "the basket cannot be modified")
	ErrBasketCannotBeCancelled        = errors.Wrap(errors.ErrBadRequest, "the basket cannot be cancelled")
	ErrQuantityCannotBeZeroOrNegative = errors.Wrap(
		errors.ErrBadRequest,
		"the item quantity cannot be zero or negative",
	)
	ErrBasketIDCannotBeBlank   = errors.Wrap(errors.ErrBadRequest, "the basket id cannot be blank")
	ErrPaymentIDCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "the payment id cannot be blank")
	ErrCustomerIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the customer id cannot be blank")
)

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*BasketES)(nil)

type BasketES struct {
	es.Aggregate
	customerID string
	paymentID  string
	items      Items
	status     BasketStatus
}

// Key implements registry.Registrable
func (BasketES) Key() string {
	return BasketAggregate
}

func (s BasketES) CustomerID() string {
	return s.customerID
}

func (s BasketES) PaymentID() string {
	return s.paymentID
}

func (s BasketES) RawItems() Items {
	return s.items
}

func (s BasketES) Items() []*Item {
	items := make([]*Item, 0, len(s.items))
	for _, item := range s.items {
		items = append(items, item)
	}

	return items
}

func (s BasketES) Status() BasketStatus {
	return s.status
}

func NewBasketES(id string) *BasketES {
	return &BasketES{
		Aggregate: es.NewAggregate(id, BasketAggregate),
		items:     make(map[string]*Item),
	}
}

func (b *BasketES) Start(customerID string) (ddd.Event, error) {
	if b.status != BasketUnknown {
		return nil, ErrBasketCannotBeModified
	}

	if customerID == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}

	b.AddEvent(BasketStartedEvent, &BasketStarted{
		CustomerID: customerID,
	})

	return ddd.NewEvent(BasketStartedEvent, b), nil
}

func (b BasketES) IsCancellable() bool {
	return b.status == BasketIsOpen
}

func (b BasketES) IsOpen() bool {
	return b.status == BasketIsOpen
}

func (b *BasketES) Cancel() (ddd.Event, error) {
	if !b.IsCancellable() {
		return nil, ErrBasketCannotBeCancelled
	}

	b.AddEvent(BasketCancelledEvent, &BasketCancelled{})

	return ddd.NewEvent(BasketCancelledEvent, b), nil
}

func (b *BasketES) Checkout(paymentID string) (ddd.Event, error) {
	if !b.IsOpen() {
		return nil, ErrBasketCannotBeModified
	}

	if len(b.items) == 0 {
		return nil, ErrBasketHasNoItems
	}

	if paymentID == "" {
		return nil, ErrPaymentIDCannotBeBlank
	}

	b.AddEvent(BasketCheckedOutEvent, &BasketCheckedOut{
		PaymentID: paymentID,
	})

	return ddd.NewEvent(BasketCheckedOutEvent, b), nil
}

func (b *BasketES) AddItem(store *Store, product *Product, quantity int) error {
	if !b.IsOpen() {
		return ErrBasketCannotBeModified
	}

	if quantity <= 0 {
		return ErrQuantityCannotBeZeroOrNegative
	}

	b.AddEvent(BasketItemAddedEvent, &BasketItemAdded{
		Item: Item{
			storeID:     store.ID,
			productID:   product.ID,
			storeName:   store.Name,
			productName: product.Name,
			price:       product.Price,
			quantity:    quantity,
		},
	})

	return nil
}

func (b *BasketES) RemoveItem(product *Product, quantity int) error {
	if !b.IsOpen() {
		return ErrBasketCannotBeModified
	}

	if quantity <= 0 {
		return ErrQuantityCannotBeZeroOrNegative
	}

	if _, exists := b.items[product.ID]; exists {
		b.AddEvent(BasketItemRemovedEvent, &BasketItemRemoved{
			ProductID: product.ID,
			Quantity:  quantity,
		})
	}

	return nil
}

func (b *BasketES) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *BasketStarted:
		b.customerID = payload.CustomerID
		b.status = BasketIsOpen

	case *BasketItemAdded:
		if item, exists := b.items[payload.Item.productID]; exists {
			item.quantity += payload.Item.quantity
			b.items[payload.Item.productID] = item
		} else {
			b.items[payload.Item.productID] = &payload.Item
		}

	case *BasketItemRemoved:
		if item, exists := b.items[payload.ProductID]; exists {
			if item.quantity-payload.Quantity <= 1 {
				delete(b.items, payload.ProductID)
			} else {
				item.quantity -= payload.Quantity
				b.items[payload.ProductID] = item
			}
		}

	case *BasketCancelled:
		b.items = make(map[string]*Item)
		b.status = BasketIsCancelled

	case *BasketCheckedOut:
		b.paymentID = payload.PaymentID
		b.status = BasketIsCheckedOut

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", b, event.EventName(), payload)
	}

	return nil
}

func (b *BasketES) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *BasketV1:
		b.customerID = ss.CustomerID
		b.paymentID = ss.PaymentID
		b.items = ss.Items
		b.status = ss.Status

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", b, snapshot)
	}

	return nil
}

func (b *BasketES) ToSnapshot() es.Snapshot {
	return &BasketV1{
		CustomerID: b.customerID,
		PaymentID:  b.paymentID,
		Items:      b.items,
		Status:     b.status,
	}
}
