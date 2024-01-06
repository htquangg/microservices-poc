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
	CreateStore struct {
		ID   string
		Name string
	}

	CreateStoreHandler decorator.CommandHandler[CreateStore]

	createStoreHandler struct {
		storeESRepo domain.StoreESRepository
		publisher   ddd.EventPublisher[ddd.Event]
		log         logger.Logger
	}
)

func NewCreateStoreHandler(
	storeESRepo domain.StoreESRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) CreateStoreHandler {
	return decorator.ApplyCommandDecorators[CreateStore](&createStoreHandler{
		storeESRepo: storeESRepo,
		publisher:   publisher,
		log:         log,
	},
		log,
	)
}

func (h *createStoreHandler) Handle(ctx context.Context, cmd CreateStore) error {
	store, err := h.storeESRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading store")
	}

	event, err := store.Init(
		cmd.Name,
	)
	if err != nil {
		return errors.Wrap(err, "initializing store")
	}

	err = h.storeESRepo.Save(ctx, store)
	if err != nil {
		return errors.Wrap(err, "error creating store")
	}

	return h.publisher.Publish(ctx, event)
}
