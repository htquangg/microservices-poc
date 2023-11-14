package command

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/decorator"
	"github.com/htquangg/microservices-poc/pkg/logger"
)

type (
	RegisterCustomer struct {
		ID    string
		Name  string
		Phone string
		Email string
	}

	RegisterCustomerHandler decorator.CommandHandler[RegisterCustomer]

	registerCustomerHandler struct {
		customerRepo domain.CustomerRepository
		publisher    ddd.EventPublisher[ddd.AggregateEvent]
		log          logger.Logger
	}
)

func NewRegisterCustomerHandler(customerRepo domain.CustomerRepository, publisher ddd.EventPublisher[ddd.AggregateEvent], log logger.Logger) RegisterCustomerHandler {
	return decorator.ApplyCommandDecorators[RegisterCustomer](
		&registerCustomerHandler{
			customerRepo: customerRepo,
			publisher:    publisher,
			log:          log,
		}, log)
}

func (h *registerCustomerHandler) Handle(ctx context.Context, cmd RegisterCustomer) error {
	customer, err := domain.RegisterCustomer(cmd.ID, cmd.Name, cmd.Phone, domain.WithCustomerEmail(cmd.Email))
	if err != nil {
		return err
	}

	if err = h.customerRepo.Save(ctx, customer); err != nil {
		return err
	}

	// publish domain events
	if err = h.publisher.Publish(ctx, customer.Events()...); err != nil {
		return err
	}

	return nil
}