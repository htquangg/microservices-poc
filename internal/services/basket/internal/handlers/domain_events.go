package handlers

import (
	"context"

	"github.com/htquangg/di/v2"
	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/basket/basketpb"
	"github.com/htquangg/microservices-poc/internal/services/basket/constants"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/domain"
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
		domain.BasketStartedEvent,
		domain.BasketCancelledEvent,
		domain.BasketCheckedOutEvent,
	)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.BasketStartedEvent:
		return h.onBasketStarted(ctx, event)
	case domain.BasketCancelledEvent:
		return h.onBasketCancelled(ctx, event)
	case domain.BasketCheckedOutEvent:
		return h.onBasketCheckedOut(ctx, event)
	}

	return nil
}

func (h domainHandlers[T]) onBasketStarted(ctx context.Context, event T) error {
	basket := event.Payload().(*domain.BasketES)
	return h.publisher.Publish(ctx,
		basketpb.BasketAggregateChannel,
		ddd.NewEvent(basketpb.BasketStartedEvent,
			&basketpb.BasketStarted{
				Id:         basket.ID(),
				CustomerId: basket.CustomerID(),
			},
		),
	)
}

func (h domainHandlers[T]) onBasketCancelled(ctx context.Context, event T) error {
	basket := event.Payload().(*domain.BasketES)
	return h.publisher.Publish(ctx,
		basketpb.BasketAggregateChannel,
		ddd.NewEvent(basketpb.BasketCancelledEvent,
			&basketpb.BasketCancelled{
				Id: basket.ID(),
			},
		),
	)
}

func (h domainHandlers[T]) onBasketCheckedOut(ctx context.Context, event T) error {
	basket := event.Payload().(*domain.BasketES)

	items := make([]*basketpb.BasketCheckedOut_Item, 0, len(basket.RawItems()))
	for _, item := range basket.RawItems() {
		items = append(items, &basketpb.BasketCheckedOut_Item{
			StoreId:     item.StoreID,
			StoreName:   item.StoreName,
			ProductId:   item.ProductID,
			ProductName: item.ProductName,
			Price:       item.Price,
			Quantity:    int32(item.Quantity),
		})
	}

	return h.publisher.Publish(ctx,
		basketpb.BasketAggregateChannel,
		ddd.NewEvent(basketpb.BasketCheckedOutEvent,
			&basketpb.BasketCheckedOut{
				Id:         basket.ID(),
				CustomerId: basket.CustomerID(),
				PaymentId:  basket.PaymentID(),
				Items:      items,
			},
		),
	)
}
