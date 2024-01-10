package commands

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/order/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/decorator"
	"github.com/htquangg/microservices-poc/pkg/logger"
	"github.com/stackus/errors"
)

type (
	CancelOrder struct {
		ID string
	}

	CancelOrderHandler decorator.CommandHandler[CancelOrder]

	cancelOrderHandler struct {
		orderESRepo domain.OrderESRepository
		publisher   ddd.EventPublisher[ddd.Event]
		log         logger.Logger
	}
)

func NewCancelOrderHandler(
	orderESRepo domain.OrderESRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) CancelOrderHandler {
	return &cancelOrderHandler{
		orderESRepo: orderESRepo,
		publisher:   publisher,
		log:         log,
	}
}

func (h *cancelOrderHandler) Handle(ctx context.Context, cmd CancelOrder) error {
	orderES, err := h.orderESRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading order")
	}

	event, err := orderES.Cancel()
	if err != nil {
		return errors.Wrap(err, "cancelling order")
	}

	if err = h.orderESRepo.Save(ctx, orderES); err != nil {
		return errors.Wrap(err, "error cancelling order")
	}

	return h.publisher.Publish(ctx, event)
}
