package decorator

import (
	"context"
	"fmt"
	"strings"

	"github.com/htquangg/microservices-poc/pkg/logger"
)

func ApplyCommandDecorators[H any](handler CommandHandler[H], log logger.Logger) CommandHandler[H] {
	return commandLoggingDecorator[H]{
		base: handler,
		log:  log,
	}
}

type CommandHandler[C any] interface {
	Handle(ctx context.Context, cmd C) error
}

func generateActionName(handler any) string {
	return strings.Split(fmt.Sprintf("%T", handler), ".")[1]
}
