package commands

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/customer/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/decorator"
	"github.com/htquangg/microservices-poc/pkg/logger"
)

type (
	AuthorizeCustomer struct {
		ID string
	}

	AuthorizeCustomerHandler decorator.CommandHandler[AuthorizeCustomer]

	authorizeCustomerHandler struct {
		customerRepo domain.CustomerRepository
		publisher    ddd.EventPublisher[ddd.AggregateEvent]
		log          logger.Logger
	}
)

func NewAuthorizeCustomerHandler(
	customerRepo domain.CustomerRepository,
	publisher ddd.EventPublisher[ddd.AggregateEvent],
	log logger.Logger,
) AuthorizeCustomerHandler {
	return decorator.ApplyCommandDecorators[AuthorizeCustomer](
		&authorizeCustomerHandler{
			customerRepo: customerRepo,
			publisher:    publisher,
			log:          log,
		},
		log,
	)
}

func (h *authorizeCustomerHandler) Handle(ctx context.Context, cmd AuthorizeCustomer) error {
	customer, err := h.customerRepo.Find(ctx, cmd.ID)
	if err != nil {
		return err
	}

	err = customer.Authorize()
	if err != nil {
		return err
	}

	return h.publisher.Publish(ctx, customer.Events()...)
}
