package handlers

import (
	"context"

	"github.com/htquangg/di/v2"
	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/store/constants"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/domain"
	pb_store "github.com/htquangg/microservices-poc/internal/services/store/proto"
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

	registerDomainEventHanders(subscriber, handlers)
}

func registerDomainEventHanders(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers, domain.ProductAddedEvent)
}

var _ ddd.EventHandler[ddd.Event] = (*domainHanlders[ddd.Event])(nil)

func (h domainHanlders[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.ProductAddedEvent:
		return h.onProductAdded(ctx, event)
	}

	return nil
}

func (h domainHanlders[T]) onProductAdded(ctx context.Context, event T) error {
	product := event.Payload().(*domain.ProductES)
	return h.publisher.Publish(
		ctx,
		pb_store.ProductAggregateChannel,
		ddd.NewEvent(pb_store.ProductAddedEvent,
			&pb_store.ProductAdded{
				Id:          product.ID(),
				StoreId:     product.StoreID(),
				Name:        product.Name(),
				Description: product.Description(),
				Sku:         product.SKU(),
				Price:       product.Price(),
			},
		),
	)
}