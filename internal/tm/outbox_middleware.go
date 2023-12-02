package tm

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/am"
)

type OutboxStore interface {
	Save(ctx context.Context, msgs ...am.Message) error
	FindUnpublished(ctx context.Context, limit int) ([]am.Message, error)
	MarkPublished(ctx context.Context, ids string) error
}

func OutboxPublisher(store OutboxStore) am.MessagePublisherMiddleware {
	return func(next am.MessagePublisher) am.MessagePublisher {
		return am.MessagePublisherFunc(func(ctx context.Context, _ string, msgs ...am.Message) error {
			err := store.Save(ctx, msgs...)
			return err
		})
	}
}
