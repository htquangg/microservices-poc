package application

import (
	"context"

	"github.com/htquangg/microservices-poc/internal/services/notification/internal/models"
)

type CustomerRepository interface {
	FindByID(ctx context.Context, id string) (*models.Customer, error)
	Add(ctx context.Context, customerID string, name string, phone string, email string) error
}
