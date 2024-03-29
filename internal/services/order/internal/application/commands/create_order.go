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
	CreateOrder struct {
		ID         string
		CustomerID string
		PaymentID  string
		Items      []*domain.Item
	}

	CreateOrderHandler decorator.CommandHandler[CreateOrder]

	createOrderHandler struct {
		orderESRepo domain.OrderESRepository
		publisher   ddd.EventPublisher[ddd.Event]
		log         logger.Logger
	}
)

func NewCreateOrderHandler(
	orderESRepo domain.OrderESRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) CreateOrderHandler {
	return &createOrderHandler{
		orderESRepo: orderESRepo,
		publisher:   publisher,
		log:         log,
	}
}

func (h *createOrderHandler) Handle(ctx context.Context, cmd CreateOrder) error {
	orderES, err := h.orderESRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading order")
	}

	event, err := orderES.CreateOrder(cmd.ID, cmd.CustomerID, cmd.PaymentID, cmd.Items)
	if err != nil {
		return errors.Wrap(err, "creating order")
	}

	if err = h.orderESRepo.Save(ctx, orderES); err != nil {
		return errors.Wrap(err, "error creating order")
	}

	return h.publisher.Publish(ctx, event)
}
