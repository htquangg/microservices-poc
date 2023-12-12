package tm

import (
	"context"
	"errors"
	"fmt"

	"github.com/htquangg/microservices-poc/internal/am"
)

type ErrDuplicateMessage string

type InboxStore interface {
	Save(ctx context.Context, msg am.IncomingMessage) error
}

func InboxHandler(store InboxStore) am.MessageHandlerMiddleware {
	return func(next am.MessageHandler) am.MessageHandler {
		return am.MessageHandlerFunc(func(ctx context.Context, msg am.IncomingMessage) error {
			if err := store.Save(ctx, msg); err != nil {
				var errDupe ErrDuplicateMessage
				if errors.As(err, &errDupe) {
					return nil
				}
				return err
			}

			return next.HandleMessage(ctx, msg)
		})
	}
}

func (e ErrDuplicateMessage) Error() string {
	return fmt.Sprintf("duplicate message id encountered: %s", string(e))
}
