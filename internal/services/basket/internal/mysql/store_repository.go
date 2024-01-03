package mysql

import (
	"context"
	"fmt"

	"github.com/htquangg/microservices-poc/internal/services/basket/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/database"
)

const StoreTable = "stores"

var _ domain.StoreRepository = (*StoreRepository)(nil)

type StoreRepository struct {
	db       database.DB
	fallback domain.StoreRepository
}

func NewStoreRepository(db database.DB, fallback domain.StoreRepository) StoreRepository {
	return StoreRepository{
		db:       db,
		fallback: fallback,
	}
}

func (r StoreRepository) FindOneByID(ctx context.Context, storeID string) (*domain.Store, error) {
	query := r.table(`
		SELECT id, name FROM %s
		WHERE id = ?
		LIMIT 1
	`)

	results, err := r.db.Engine(ctx).Query(query, storeID)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, err
	}

	return &domain.Store{
		ID:   string(results[0]["id"]),
		Name: string(results[0]["name"]),
	}, nil
}

func (StoreRepository) table(query string) string {
	return fmt.Sprintf(query, StoreTable)
}
