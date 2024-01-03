package domain

import (
	"context"
)

type StoreRepository interface {
	FindOneByID(ctx context.Context, storeID string) (*Store, error)
}
