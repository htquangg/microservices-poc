package decorator

import (
	"context"

	"github.com/htquangg/microservices-poc/pkg/logger"
)

func ApplyQueryDecorators[H any, R any](handler QueryHandler[H, R], log logger.Logger) QueryHandler[H, R] {
	return queryLoggingDecorator[H, R]{
		base: handler,
		log:  log,
	}
}

type QueryHandler[Q any, R any] interface {
	Handle(ctx context.Context, q Q) (R, error)
}
