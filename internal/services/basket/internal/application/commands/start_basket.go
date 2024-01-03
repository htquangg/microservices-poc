package commands

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/basket/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/decorator"
	"github.com/htquangg/microservices-poc/pkg/logger"

	"github.com/stackus/errors"
)

type (
	StartBasket struct {
		ID         string
		CustomerID string
	}

	StartBasketHandler decorator.CommandHandler[StartBasket]

	startBasketHandler struct {
		basketESRepo domain.BasketESRepository
		publisher    ddd.EventPublisher[ddd.Event]
		log          logger.Logger
	}
)

func NewStartBasketHandler(
	basketESRepo domain.BasketESRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) StartBasketHandler {
	return &startBasketHandler{
		basketESRepo: basketESRepo,
		publisher:    publisher,
		log:          log,
	}
}

func (h *startBasketHandler) Handle(ctx context.Context, cmd StartBasket) error {
	basketES, err := h.basketESRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading basket")
	}

	event, err := basketES.Start(cmd.CustomerID)
	if err != nil {
		return errors.Wrap(err, "initializing basket")
	}

	if err = h.basketESRepo.Save(ctx, basketES); err != nil {
		return errors.Wrap(err, "error creating basket")
	}

	return h.publisher.Publish(ctx, event)
}
