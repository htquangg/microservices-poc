package domain

import (
	"context"
)

type StoreRepository interface {
	Add(ctx context.Context, store *Store) error
	Rebrand(ctx context.Context, id string, name string) error
	FindOneByID(ctx context.Context, storeID string) (*Store, error)
}
