package handlers

import (
	"context"

	"github.com/htquangg/di/v2"
	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/order/constants"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/domain"
	"github.com/htquangg/microservices-poc/internal/services/order/orderpb"
)

type domainHandlers[T ddd.Event] struct {
	publisher am.EventPublisher
}

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return domainHandlers[ddd.Event]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		domainHandlers := di.Get(ctx, constants.DomainEventHandlersKey).(ddd.EventHandler[ddd.Event])

		return domainHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	registerDomainEventHandlers(subscriber, handlers)
}

func registerDomainEventHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.OrderCreatedEvent,
		domain.OrderRejectedEvent,
		domain.OrderApprovedEvent,
		domain.OrderCancelledEvent,
		domain.OrderCompletedEvent,
	)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.OrderCreatedEvent:
		return h.onOrderCreated(ctx, event)
	case domain.OrderRejectedEvent:
		return h.onOrderRejected(ctx, event)
	case domain.OrderApprovedEvent:
		return h.onOrderApproved(ctx, event)
	case domain.OrderCancelledEvent:
		return h.onOrderCancelled(ctx, event)
	case domain.OrderCompletedEvent:
		return h.onOrderCompleted(ctx, event)
	}

	return nil
}

func (h domainHandlers[T]) onOrderCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.OrderES)
	items := make([]*orderpb.OrderCreated_Item, 0, len(payload.Items()))
	for _, item := range payload.Items() {
		items = append(items, &orderpb.OrderCreated_Item{
			ProductId: item.ProductID(),
			StoreId:   item.StoreID(),
			Price:     item.Price(),
			Quantity:  int32(item.Quantity()),
		})
	}

	return h.publisher.Publish(ctx,
		orderpb.OrderAggregateChannel,
		ddd.NewEvent(orderpb.OrderCreatedEvent,
			&orderpb.OrderCreated{
				Id:         payload.ID(),
				CustomerId: payload.CustomerID(),
				PaymentId:  payload.PaymentID(),
				ShoppingId: payload.ShoppingID(),
				Items:      items,
			},
		),
	)
}

func (h domainHandlers[T]) onOrderRejected(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.OrderES)
	return h.publisher.Publish(ctx,
		orderpb.OrderAggregateChannel,
		ddd.NewEvent(orderpb.OrderRejectedEvent,
			&orderpb.OrderRejected{
				Id:         payload.ID(),
				CustomerId: payload.CustomerID(),
				PaymentId:  payload.PaymentID(),
			},
		),
	)
}

func (h domainHandlers[T]) onOrderApproved(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.OrderES)
	return h.publisher.Publish(ctx,
		orderpb.OrderAggregateChannel,
		ddd.NewEvent(orderpb.OrderApprovedEvent,
			&orderpb.OrderApproved{
				Id:         payload.ID(),
				CustomerId: payload.CustomerID(),
				PaymentId:  payload.PaymentID(),
			},
		),
	)
}

func (h domainHandlers[T]) onOrderCancelled(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.OrderES)
	return h.publisher.Publish(ctx,
		orderpb.OrderAggregateChannel,
		ddd.NewEvent(orderpb.OrderCanceledEvent,
			&orderpb.OrderCanceled{
				Id:         payload.ID(),
				CustomerId: payload.CustomerID(),
				PaymentId:  payload.CustomerID(),
			},
		),
	)
}

func (h domainHandlers[T]) onOrderCompleted(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.OrderES)
	return h.publisher.Publish(ctx,
		orderpb.OrderAggregateChannel,
		ddd.NewEvent(orderpb.OrderCompletedEvent,
			&orderpb.OrderCompleted{
				Id:         payload.ID(),
				CustomerId: payload.CustomerID(),
				InvoiceId:  payload.InvoiceID(),
			},
		),
	)
}
