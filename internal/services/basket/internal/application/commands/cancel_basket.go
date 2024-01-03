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
	CancelBasket struct {
		ID string
	}

	CancelBasketHandler decorator.CommandHandler[CancelBasket]

	cancelBasketHandler struct {
		basketESRepo domain.BasketESRepository
		publisher    ddd.EventPublisher[ddd.Event]
		log          logger.Logger
	}
)

func NewCancelBasketHandler(
	basketESRepo domain.BasketESRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) CancelBasketHandler {
	return &cancelBasketHandler{
		basketESRepo: basketESRepo,
		publisher:    publisher,
		log:          log,
	}
}

func (h *cancelBasketHandler) Handle(ctx context.Context, cmd CancelBasket) error {
	basketES, err := h.basketESRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading basket")
	}

	event, err := basketES.Cancel()
	if err != nil {
		return errors.Wrap(err, "cancelling basket")
	}

	if err = h.basketESRepo.Save(ctx, basketES); err != nil {
		return errors.Wrap(err, "error cancelling basket")
	}

	return h.publisher.Publish(ctx, event)
}
