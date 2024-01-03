package domain

import (
	"context"
)

type ProductRepository interface {
	FindOneByID(ctx context.Context, productID string) (*Product, error)
}
