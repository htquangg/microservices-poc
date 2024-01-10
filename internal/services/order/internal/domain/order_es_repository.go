package domain

import "context"

type OrderESRepository interface {
	Load(ctx context.Context, id string) (*OrderES, error)
	Save(ctx context.Context, order *OrderES) error
}
