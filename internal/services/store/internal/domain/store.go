package domain

import "context"

type Store struct {
	ID   string
	Name string
}

type StoreRepository interface {
	AddStore(ctx context.Context, id string, name string) error
	RenameStore(ctx context.Context, id string, name string) error
	FindOneByID(ctx context.Context, id string) (*Store, error)
	FindAll(ctx context.Context) ([]*Store, error)
}
