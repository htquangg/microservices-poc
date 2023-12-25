package handlers

import (
	"context"

	"github.com/htquangg/di/v2"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/store/constants"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/domain"
)

type storeHandlers[T ddd.Event] struct {
	storeRepo domain.StoreRepository
}

var _ ddd.EventHandler[ddd.Event] = (*storeHandlers[ddd.Event])(nil)

func NewStoreHandlers(storeRepo domain.StoreRepository) ddd.EventHandler[ddd.Event] {
	return storeHandlers[ddd.Event]{
		storeRepo: storeRepo,
	}
}

func RegisterStoreHandlers(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.Event](func(ctx context.Context, event ddd.Event) error {
		storeHandlers := di.Get(ctx, constants.StoreHandlersKey).(ddd.EventHandler[ddd.Event])

		return storeHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.Event])

	registerStoreHandlers(subscriber, handlers)
}

func registerStoreHandlers(subscriber ddd.EventSubscriber[ddd.Event], handlers ddd.EventHandler[ddd.Event]) {
	subscriber.Subscribe(handlers,
		domain.StoreCreatedEvent,
		domain.StoreRebrandedEvent,
	)
}

func (h storeHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.StoreCreatedEvent:
		h.onStoreCreated(ctx, event)
	case domain.StoreRebrandedEvent:
		h.onStoreRebranded(ctx, event)

	}

	return nil
}

func (h storeHandlers[T]) onStoreCreated(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.StoreES)
	return h.storeRepo.AddStore(ctx, payload.ID(), payload.Name())
}

func (h storeHandlers[T]) onStoreRebranded(ctx context.Context, event T) error {
	payload := event.Payload().(*domain.StoreES)
	return h.storeRepo.RenameStore(ctx, payload.ID(), payload.Name())
}
