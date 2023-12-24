package domain

import (
	"context"
)

type ProductESRepository interface {
	Load(ctx context.Context, id string) (*ProductES, error)
	Save(ctx context.Context, product *ProductES) error
}
