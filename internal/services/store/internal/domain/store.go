package domain

import "context"

type Store struct {
	ID   string
	Name string
}

type StoreRepository interface {
	AddStore(ctx context.Context, id string, name string) error
	RenameStore(ctx context.Context, id string, name string) error
	One(ctx context.Context, id string) (*Store, error)
	All(ctx context.Context) ([]*Store, error)
}
