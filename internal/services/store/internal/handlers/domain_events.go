package handlers

import (
	"context"

	"github.com/htquangg/di/v2"
	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/store/constants"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/domain"
	"github.com/htquangg/microservices-poc/internal/services/store/storepb"
)

type domainHanlders[T ddd.Event] struct {
	publisher am.EventPublisher
}

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.Event] {
	return &domainHanlders[ddd.Event]{
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
		domain.StoreCreatedEvent,
		domain.StoreRebrandedEvent,
		domain.ProductAddedEvent,
	)
}

var _ ddd.EventHandler[ddd.Event] = (*domainHanlders[ddd.Event])(nil)

func (h domainHanlders[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.StoreCreatedEvent:
		return h.onStoreCreated(ctx, event)
	case domain.StoreRebrandedEvent:
		return h.onStoreRebranded(ctx, event)
	case domain.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	}

	return nil
}

func (h domainHanlders[T]) onStoreCreated(ctx context.Context, event T) error {
	product := event.Payload().(*domain.StoreES)
	return h.publisher.Publish(
		ctx,
		storepb.StoreAggregateChannel,
		ddd.NewEvent(storepb.StoreCreatedEvent,
			&storepb.StoreCreated{
				Id:   product.ID(),
				Name: product.Name(),
			},
		),
	)
}

func (h domainHanlders[T]) onStoreRebranded(ctx context.Context, event T) error {
	product := event.Payload().(*domain.StoreES)
	return h.publisher.Publish(
		ctx,
		storepb.StoreAggregateChannel,
		ddd.NewEvent(storepb.StoreRebrandedEvent,
			&storepb.StoreRebranded{
				Id:   product.ID(),
				Name: product.Name(),
			},
		),
	)
}

func (h domainHanlders[T]) onProductAdded(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.ProductES)
	return h.publisher.Publish(
		ctx,
		storepb.ProductAggregateChannel,
		ddd.NewEvent(storepb.ProductAddedEvent,
			&storepb.ProductAdded{
				Id:          payload.ID(),
				StoreId:     payload.StoreID(),
				Name:        payload.Name(),
				Description: payload.Description(),
				Sku:         payload.SKU(),
				Price:       payload.Price(),
			},
		),
	)
}
