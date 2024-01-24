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
	ApproveOrder struct {
		ID         string
		ShoppingID string
	}

	ApproveOrderHandler decorator.CommandHandler[ApproveOrder]

	approveOrderHander struct {
		orderESRepo domain.OrderESRepository
		publisher   ddd.EventPublisher[ddd.Event]
		log         logger.Logger
	}
)

func NewApproveOrderHandler(
	orderESRepo domain.OrderESRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) ApproveOrderHandler {
	return &approveOrderHander{
		orderESRepo: orderESRepo,
		publisher:   publisher,
		log:         log,
	}
}

func (h *approveOrderHander) Handle(ctx context.Context, cmd ApproveOrder) error {
	orderES, err := h.orderESRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading order")
	}

	event, err := orderES.Approve(cmd.ShoppingID)
	if err != nil {
		return errors.Wrap(err, "approving order")
	}

	if err = h.orderESRepo.Save(ctx, orderES); err != nil {
		return errors.Wrap(err, "error approving order")
	}

	return h.publisher.Publish(ctx, event)
}
