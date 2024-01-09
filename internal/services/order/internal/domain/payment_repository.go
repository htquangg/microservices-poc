package domain

import "context"

type PaymentRepository interface {
	Confirm(ctx context.Context, id string) error
}
