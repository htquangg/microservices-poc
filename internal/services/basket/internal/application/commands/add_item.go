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
	AddItem struct {
		ID        string
		ProductID string
		Quantity  int
	}

	AddItemHandler decorator.CommandHandler[AddItem]

	addItemHandler struct {
		basketESRepo domain.BasketESRepository
		storeRepo    domain.StoreRepository
		productRepo  domain.ProductRepository
		publisher    ddd.EventPublisher[ddd.Event]
		log          logger.Logger
	}
)

func NewAddItemHandler(
	basketESRepo domain.BasketESRepository,
	storeRepo domain.StoreRepository,
	productRepo domain.ProductRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) AddItemHandler {
	return &addItemHandler{
		basketESRepo: basketESRepo,
		storeRepo:    storeRepo,
		productRepo:  productRepo,
		publisher:    publisher,
		log:          log,
	}
}

func (h *addItemHandler) Handle(ctx context.Context, cmd AddItem) error {
	basketES, err := h.basketESRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading basket")
	}

	product, err := h.productRepo.FindOneByID(ctx, cmd.ProductID)
	if err != nil {
		return err
	}

	store, err := h.storeRepo.FindOneByID(ctx, product.StoreID)
	if err != nil {
		return err
	}

	err = basketES.AddItem(store, product, cmd.Quantity)
	if err != nil {
		return errors.Wrap(err, "adding item into basket")
	}

	if err = h.basketESRepo.Save(ctx, basketES); err != nil {
		return errors.Wrap(err, "error adding item into basket")
	}

	return nil
}
