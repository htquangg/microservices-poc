package domain

import "context"

type StoreESRepository interface {
	Load(ctx context.Context, id string) (*StoreES, error)
	Save(ctx context.Context, store *StoreES) error
}
