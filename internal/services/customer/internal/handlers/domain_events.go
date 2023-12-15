package handlers

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/customer/constants"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/domain"
	pb_customer "github.com/htquangg/microservices-poc/internal/services/customer/proto"

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

func RegisterDomainEventHandlers(container di.Container) {
	handlers := ddd.EventHandlerFunc[ddd.AggregateEvent](func(ctx context.Context, event ddd.AggregateEvent) error {
		domainHandlers := di.Get(ctx, (constants.DomainEventHandlersKey)).(ddd.EventHandler[ddd.AggregateEvent])

		return domainHandlers.HandleEvent(ctx, event)
	})

	subscriber := container.Get(constants.DomainDispatcherKey).(*ddd.EventDispatcher[ddd.AggregateEvent])

	registerDomainEventHandlers(subscriber, handlers)
}

func registerDomainEventHandlers(
	subscriber ddd.EventSubscriber[ddd.AggregateEvent],
	handlers ddd.EventHandler[ddd.AggregateEvent],
) {
	subscriber.Subscribe(handlers, domain.CustomerRegisteredEvent)
}

func (h domainHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case domain.CustomerRegisteredEvent:
		return h.onCustomerRegistered(ctx, event)
	}
	return nil
}

func (h domainHandlers[T]) onCustomerRegistered(ctx context.Context, event ddd.AggregateEvent) error {
	payload := event.Payload().(*domain.CustomerRegistered)
	return h.publisher.Publish(
		ctx,
		pb_customer.CustomerAggregateChannel,
		ddd.NewEvent(pb_customer.CustomerRegisteredEvent, &pb_customer.CustomerRegistered{
			Id:    payload.Customer.ID(),
			Name:  payload.Customer.Name(),
			Phone: payload.Customer.Phone(),
			Email: payload.Customer.Email(),
		}),
	)
}
