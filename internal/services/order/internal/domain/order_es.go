package domain

import (
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/es"

	"github.com/stackus/errors"
)

const OrderAggregate = "orders.Order"

var (
	ErrOrderAlreadyCreated     = errors.Wrap(errors.ErrBadRequest, "the order cannot be recreated")
	ErrOrderHasNoItems         = errors.Wrap(errors.ErrBadRequest, "the order has no items")
	ErrOrderCannotBeCancelled  = errors.Wrap(errors.ErrBadRequest, "the order cannot be cancelled")
	ErrCustomerIDCannotBeBlank = errors.Wrap(errors.ErrBadRequest, "the customer id cannot be blank")
	ErrPaymentIDCannotBeBlank  = errors.Wrap(errors.ErrBadRequest, "the payment id cannot be blank")
)

type OrderES struct {
	es.Aggregate
	customerID string
	paymentID  string
	invoiceID  string
	shoppingID string
	status     OrderStatus
	items      []*Item
}

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*OrderES)(nil)

func NewOrderES(id string) *OrderES {
	return &OrderES{
		Aggregate: es.NewAggregate(id, OrderAggregate),
	}
}

// Key implements registry.Registrable
func (OrderES) Key() string {
	return OrderAggregate
}

func (o OrderES) CustomerID() string {
	return o.customerID
}

func (o OrderES) PaymentID() string {
	return o.paymentID
}

func (o OrderES) InvoiceID() string {
	return o.invoiceID
}

func (o OrderES) ShoppingID() string {
	return o.shoppingID
}

func (o OrderES) Status() OrderStatus {
	return o.status
}

func (o OrderES) Items() []*Item {
	return o.items
}

func (o *OrderES) CreateOrder(id string, customerID string, paymentID string, items []*Item) (ddd.Event, error) {
	if o.status != OrderUnknown {
		return nil, ErrOrderAlreadyCreated
	}

	if len(items) == 0 {
		return nil, ErrOrderHasNoItems
	}

	if customerID == "" {
		return nil, ErrCustomerIDCannotBeBlank
	}

	if paymentID == "" {
		return nil, ErrPaymentIDCannotBeBlank
	}

	o.AddEvent(OrderCreatedEvent, &OrderCreated{
		CustomerID: customerID,
		PaymentID:  paymentID,
		Items:      items,
	})

	return ddd.NewEvent(OrderCreatedEvent, o), nil
}

func (o *OrderES) Reject() (ddd.Event, error) {
	o.AddEvent(OrderRejectedEvent, &OrderRejected{})

	return ddd.NewEvent(OrderRejectedEvent, o), nil
}

func (o *OrderES) Approve(shoppingID string) (ddd.Event, error) {
	o.AddEvent(OrderApprovedEvent, &OrderApproved{
		ShoppingID: shoppingID,
	})

	return ddd.NewEvent(OrderApprovedEvent, o), nil
}

func (o *OrderES) Cancel() (ddd.Event, error) {
	if o.status != OrderIsPending {
		return nil, ErrOrderCannotBeCancelled
	}

	o.AddEvent(OrderCancelledEvent, &OrderCancelled{
		CustomerID: o.customerID,
		PaymentID:  o.paymentID,
	})
	return ddd.NewEvent(OrderCancelledEvent, o), nil
}

func (o *OrderES) Complete(invoiceID string) (ddd.Event, error) {
	o.AddEvent(OrderCompletedEvent, &OrderCompleted{
		CustomerID: o.customerID,
		InvoiceID:  invoiceID,
	})

	return ddd.NewEvent(OrderCompletedEvent, o), nil
}

func (o *OrderES) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *OrderCreated:
		o.customerID = payload.CustomerID
		o.paymentID = payload.PaymentID
		o.items = payload.Items
		o.status = OrderIsPending

	case *OrderRejected:
		o.status = OrderIsRejected

	case *OrderApproved:
		o.shoppingID = payload.ShoppingID
		o.status = OrderIsApproved

	case *OrderCancelled:
		o.status = OrderIsCancelled

	case *OrderCompleted:
		o.invoiceID = payload.InvoiceID
		o.status = OrderIsCompleted

	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", o, event.EventName(), payload)
	}

	return nil
}

func (o *OrderES) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *OrderV1:
		o.customerID = ss.CustomerID
		o.paymentID = ss.PaymentID
		o.invoiceID = ss.InvoiceID
		o.shoppingID = ss.ShoppingID
		o.items = ss.Items
		o.status = ss.Status

	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", o, snapshot)
	}

	return nil
}

func (o *OrderES) ToSnapshot() es.Snapshot {
	return &OrderV1{
		CustomerID: o.customerID,
		PaymentID:  o.paymentID,
		InvoiceID:  o.invoiceID,
		ShoppingID: o.shoppingID,
		Items:      o.items,
		Status:     o.status,
	}
}
