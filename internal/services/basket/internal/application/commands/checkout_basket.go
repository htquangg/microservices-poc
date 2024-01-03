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
	CheckoutBasket struct {
		ID        string
		PaymentID string
	}

	CheckoutBasketHandler decorator.CommandHandler[CheckoutBasket]

	checkoutBasketHandler struct {
		basketESRepo domain.BasketESRepository
		publisher    ddd.EventPublisher[ddd.Event]
		log          logger.Logger
	}
)

func NewCheckoutBasketHandler(
	basketESRepo domain.BasketESRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) CheckoutBasketHandler {
	return &checkoutBasketHandler{
		basketESRepo: basketESRepo,
		publisher:    publisher,
		log:          log,
	}
}

func (h *checkoutBasketHandler) Handle(ctx context.Context, cmd CheckoutBasket) error {
	basketES, err := h.basketESRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading basket")
	}

	event, err := basketES.Checkout(cmd.PaymentID)
	if err != nil {
		return errors.Wrap(err, "basket checkout")
	}

	if err = h.basketESRepo.Save(ctx, basketES); err != nil {
		return errors.Wrap(err, "error basket checkout")
	}

	return h.publisher.Publish(ctx, event)
}
