package mysql

import (
	"context"
	"fmt"

	"github.com/htquangg/microservices-poc/internal/services/basket/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/converter"
	"github.com/htquangg/microservices-poc/pkg/database"
)

const ProductTable = "products"

var _ domain.ProductRepository = (*ProductRepository)(nil)

type ProductRepository struct {
	db       database.DB
	fallback domain.ProductRepository
}

func NewProductRepository(db database.DB, fallback domain.ProductRepository) ProductRepository {
	return ProductRepository{
		db:       db,
		fallback: fallback,
	}
}

func (r ProductRepository) FindOneByID(ctx context.Context, productID string) (*domain.Product, error) {
	query := r.table(`
		SELECT id, store_id, name, price FROM %s
		WHERE id = ?
		LIMIT 1
	`)

	results, err := r.db.Engine(ctx).Query(query, productID)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, err
	}

	return &domain.Product{
		ID:      string(results[0]["id"]),
		StoreID: string(results[0]["store_id"]),
		Name:    string(results[0]["name"]),
		Price:   converter.StringToFloat64(string(results[0]["price"])),
	}, nil
}

func (r ProductRepository) table(query string) string {
	return fmt.Sprintf(query, ProductTable)
}
