package handlers

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/customer/constants"
	"github.com/htquangg/microservices-poc/internal/services/customer/customerpb"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/domain"

	"github.com/htquangg/di/v2"
)

type domainHandlers[T ddd.AggregateEvent] struct {
	publisher am.EventPublisher
}

var _ ddd.EventHandler[ddd.AggregateEvent] = (*domainHandlers[ddd.AggregateEvent])(nil)

func NewDomainEventHandlers(publisher am.EventPublisher) ddd.EventHandler[ddd.AggregateEvent] {
	return &domainHandlers[ddd.AggregateEvent]{
		publisher: publisher,
	}
}

func RegisterDomainEventHandlers(ctn di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.AggregateEvent](func(ctx context.Context, event ddd.AggregateEvent) error {
		domainHandlers := ctn.Get(constants.DomainEventHandlersKey).(ddd.EventHandler[ddd.AggregateEvent])

		return domainHandlers.HandleEvent(ctx, event)
	})

	subscriber := ctn.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.AggregateEvent])

	registerDomainEventHandlers(subscriber, handlers)
}

func registerDomainEventHandlers(
	subscriber ddd.EventSubscriber[ddd.AggregateEvent],
	handlers ddd.EventHandler[ddd.AggregateEvent],
) {
	subscriber.Subscribe(handlers, domain.CustomerRegisteredEvent)
}

func (h *domainHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.CustomerRegisteredEvent:
		return h.onCustomerRegistered(ctx, event)
	}
	return nil
}

func (h *domainHandlers[T]) onCustomerRegistered(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.CustomerRegistered)
	return h.publisher.Publish(
		ctx,
		customerpb.CustomerAggregateChannel,
		ddd.NewEvent(customerpb.CustomerRegisteredEvent,
			&customerpb.CustomerRegistered{
				Id:    payload.Customer.ID(),
				Name:  payload.Customer.Name(),
				Phone: payload.Customer.Phone(),
				Email: payload.Customer.Email(),
			}),
	)
}
