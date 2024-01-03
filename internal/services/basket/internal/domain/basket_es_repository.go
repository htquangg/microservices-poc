package domain

import (
	"context"
)

type BasketESRepository interface {
	Load(ctx context.Context, basketID string) (*BasketES, error)
	Save(ctx context.Context, basket *BasketES) error
}
