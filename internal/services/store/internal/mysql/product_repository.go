package mysql

import (
	"context"
	"fmt"

	"github.com/htquangg/microservices-poc/internal/services/store/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/database"
)

const ProductTable = "products"

type ProductRepository struct {
	db database.DB
}

var _ domain.ProductRepository = (*ProductRepository)(nil)

func NewProductRepository(db database.DB) ProductRepository {
	return ProductRepository{
		db: db,
	}
}

func (r ProductRepository) AddProduct(
	ctx context.Context,
	id string,
	storeID string,
	name string,
	description string,
	sku string,
	price float64,
) error {
	query := r.table(`
		INSERT INTO %s (id, store_id, name, description, sku, price)
		VALUES (?, ?, ?, ?, ?, ?)
	`)

	_, err := r.db.Exec(ctx, query, id, storeID, name, description, sku, price)

	return err
}

func (r ProductRepository) table(query string) string {
	return fmt.Sprintf(query, ProductTable)
}
