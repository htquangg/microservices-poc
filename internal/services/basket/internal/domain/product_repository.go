package domain

import (
	"context"
)

type ProductRepository interface {
	Add(ctx context.Context, productID string, storeID string, name string, sku string, price float64) error
	FindOneByID(ctx context.Context, productID string) (*Product, error)
}
