package handlers

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/services/customer/customerpb"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/application"
)

type integrationHandlers[T ddd.Event] struct {
	app          *application.Application
	customerRepo application.CustomerRepository
}

var _ ddd.EventHandler[ddd.Event] = (*integrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(
	reg registry.Registry,
	app *application.Application,
	customerRepo application.CustomerRepository,
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewEventHandler(reg, integrationHandlers[ddd.Event]{
		app:          app,
		customerRepo: customerRepo,
	}, mws...)
}

func RegisterIntegrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(customerpb.CustomerAggregateChannel, handlers, am.MessageFilter{
		customerpb.CustomerRegisteredEvent,
	})
	return err
}

func (h integrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case customerpb.CustomerRegisteredEvent:
		return h.onCustomerRegistered(ctx, event)
	}

	return nil
}

func (h integrationHandlers[T]) onCustomerRegistered(ctx context.Context, event T) error {
	payload := (event.Payload()).(*customerpb.CustomerRegistered)
	return h.customerRepo.Add(ctx, payload.GetId(), payload.GetName(), payload.GetPhone(), payload.GetEmail())
}
