package handlers

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/store/constants"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/domain"

	"github.com/htquangg/di/v2"
)

type productHandlers[T ddd.Event] struct {
	productRepo domain.ProductRepository
}

var _ ddd.EventHandler[ddd.Event] = (*productHandlers[ddd.Event])(nil)

func NewProductHandlers(productRepo domain.ProductRepository) ddd.EventHandler[ddd.Event] {
	return productHandlers[ddd.Event]{
		productRepo: productRepo,
	}
}

func RegisterProductHandlers(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		productHandlers := di.Get(ctx, constants.ProductHandlersKey).(ddd.EventHandler[ddd.Event])

		return productHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	registerProductHandlers(subscriber, handlers)
}

func registerProductHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers, domain.ProductAddedEvent)
}

func (h productHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.ProductAddedEvent:
		h.onProductAdded(ctx, event)
	}

	return nil
}

func (h productHandlers[T]) onProductAdded(ctx context.Context, event ddd.Event) error {
	product := event.Payload().(*domain.ProductES)
	return h.productRepo.AddProduct(
		ctx,
		product.ID(),
		product.StoreID(),
		product.Name(),
		product.Description(),
		product.SKU(),
		product.Price(),
	)
}
