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
	RebrandStore struct {
		ID   string
		Name string
	}

	RebrandStoreHandler decorator.CommandHandler[RebrandStore]

	rebrandStoreHandler struct {
		storeESRepo domain.StoreESRepository
		publisher   ddd.EventPublisher[ddd.Event]
		log         logger.Logger
	}
)

func NewRebrandStoreHandler(
	storeESRepo domain.StoreESRepository,
	publisher ddd.EventPublisher[ddd.Event],
	log logger.Logger,
) RebrandStoreHandler {
	return decorator.ApplyCommandDecorators[RebrandStore](&rebrandStoreHandler{
		storeESRepo: storeESRepo,
		publisher:   publisher,
		log:         log,
	},
		log,
	)
}

func (h *rebrandStoreHandler) Handle(ctx context.Context, cmd RebrandStore) error {
	store, err := h.storeESRepo.Load(ctx, cmd.ID)
	if err != nil {
		return errors.Wrap(err, "error loading store")
	}

	event, err := store.Rebrand(
		cmd.Name,
	)
	if err != nil {
		return errors.Wrap(err, "rebranding store")
	}

	err = h.storeESRepo.Save(ctx, store)
	if err != nil {
		return errors.Wrap(err, "error rebranding store")
	}

	return h.publisher.Publish(ctx, event)
}
