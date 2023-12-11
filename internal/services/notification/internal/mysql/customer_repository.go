package mysql

import (
	"context"
	"fmt"

	"github.com/htquangg/microservices-poc/internal/services/notification/internal/application"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/constants"
	"github.com/htquangg/microservices-poc/internal/services/notification/internal/models"
	"github.com/htquangg/microservices-poc/pkg/database"
)

type CustomerRepository struct {
	db *database.DB
}

var _ application.CustomerRepository = (*CustomerRepository)(nil)

func NewCustomerRepository(db *database.DB) CustomerRepository {
	return CustomerRepository{
		db: db,
	}
}

func (r CustomerRepository) Add(ctx context.Context, customerID string, name string, phone string, email string) error {
	query := fmt.Sprintf(
		"INSERT INTO %s (id, name, phone, email) VALUES (?, ?, ?, ?) ON DUPLICATE KEY UPDATE id=id",
		constants.CustomerTableName,
	)

	_, err := r.db.Exec(ctx, query, customerID, name, phone, email)

	return err
}

func (CustomerRepository) FindByID(ctx context.Context, id string) (*models.Customer, error) {
	panic("unimplemented")
}
