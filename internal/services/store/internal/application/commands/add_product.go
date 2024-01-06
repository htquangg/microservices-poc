package commands

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/services/store/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/decorator"
	"github.com/htquangg/microservices-poc/pkg/logger"

	"github.com/stackus/errors"
)

type (
	AddProduct struct {
		ID          string
		StoreID     string
		Name        string
		Description string
		SKU         string
		Price       float64
	}

	AddProductHandler decorator.CommandHandler[AddProduct]

	addProductHandler struct {
		productESRepo domain.ProductESRepository
		publisher     ddd.EventPublisher[ddd.Event]
		log           logger.Logger
	}
)

func NewAddProductHandler(
	productESRepo domain.ProductESRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) AddProductHandler {
	return decorator.ApplyCommandDecorators[AddProduct](
		&addProductHandler{
			productESRepo: productESRepo,
			publisher:     publisher,
			log:           log,
		},
		log,
	)
}

func (h *addProductHandler) Handle(ctx context.Context, cmd AddProduct) error {
	product, err := h.productESRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading product")
	}

	event, err := product.Init(
		cmd.ID,
		cmd.StoreID,
		cmd.Name,
		cmd.Description,
		cmd.SKU,
		cmd.Price,
	)
	if err != nil {
		return errors.Wrap(err, "initializing product")
	}

	err = h.productESRepo.Save(ctx, product)
	if err != nil {
		return errors.Wrap(err, "error adding product")
	}

	return errors.Wrap(h.publisher.Publish(ctx, event), "publishing domain event")
}
