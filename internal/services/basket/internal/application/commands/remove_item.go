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
	RemoveItem struct {
		ID        string
		ProductID string
		Quantity  int
	}

	RemoveItemHandler decorator.CommandHandler[RemoveItem]

	removeItemHandler struct {
		basketESRepo domain.BasketESRepository
		productRepo  domain.ProductRepository
		publisher    ddd.EventPublisher[ddd.Event]
		log          logger.Logger
	}
)

func NewRemoveItemHandler(
	basketESRepo domain.BasketESRepository,
	productRepo domain.ProductRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) RemoveItemHandler {
	return &removeItemHandler{
		basketESRepo: basketESRepo,
		productRepo:  productRepo,
		publisher:    publisher,
		log:          log,
	}
}

func (h *removeItemHandler) Handle(ctx context.Context, cmd RemoveItem) error {
	product, err := h.productRepo.FindOneByID(ctx, cmd.ProductID)
	if err != nil {
		return err
	}

	basketES, err := h.basketESRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading basket")
	}

	err = basketES.RemoveItem(product, cmd.Quantity)
	if err != nil {
		return errors.Wrap(err, "removing item from basket")
	}

	if err = h.basketESRepo.Save(ctx, basketES); err != nil {
		return errors.Wrap(err, "error removing item from basket")
	}

	return nil
}
