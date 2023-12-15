package handlers

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/application"

	pb_customer "github.com/htquangg/microservices-poc/internal/services/customer/proto"
)

type intergrationHandlers[T ddd.Event] struct {
	app          *application.Application
	customerRepo application.CustomerRepository
}

var _ ddd.EventHandler[ddd.Event] = (*intergrationHandlers[ddd.Event])(nil)

func NewIntegrationEventHandlers(
	reg registry.Registry,
	app *application.Application,
	customerRepo application.CustomerRepository,
	mws ...am.MessageHandlerMiddleware,
) am.MessageHandler {
	return am.NewEventHandler(reg, intergrationHandlers[ddd.Event]{
		app:          app,
		customerRepo: customerRepo,
	}, mws...)
}

func RegisterIntergrationEventHandlers(subscriber am.MessageSubscriber, handlers am.MessageHandler) error {
	_, err := subscriber.Subscribe(pb_customer.CustomerAggregateChannel, handlers, am.MessageFilter{
		pb_customer.CustomerRegisteredEvent,
	})
	return err
}

func (h intergrationHandlers[T]) HandleEvent(ctx context.Context, event T) error {
	switch event.EventName() {
	case pb_customer.CustomerRegisteredEvent:
		return h.onCustomerRegistered(ctx, event)
	}

	return nil
}

func (h intergrationHandlers[T]) onCustomerRegistered(ctx context.Context, event T) error {
	payload := (event.Payload()).(*pb_customer.CustomerRegistered)
	return h.customerRepo.Add(ctx, payload.GetId(), payload.GetName(), payload.GetPhone(), payload.GetEmail())
}
