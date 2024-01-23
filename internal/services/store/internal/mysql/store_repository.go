package mysql

import (
	"context"
	"fmt"

	"github.com/htquangg/microservices-poc/internal/services/store/internal/domain"
	"github.com/htquangg/microservices-poc/pkg/database"
)

const StoreTable = "stores"

type StoreRepository struct {
	db database.DB
}

var _ domain.StoreRepository = (*StoreRepository)(nil)

func NewStoreRepository(db database.DB) *StoreRepository {
	return &StoreRepository{
		db: db,
	}
}

func (r *StoreRepository) AddStore(ctx context.Context, id string, name string) error {
	query := r.table(`
		INSERT INTO %s (id, name)
		VALUES (?, ?)
	`)

	_, err := r.db.Exec(ctx, query, id, name)

	return err
}

func (r *StoreRepository) RenameStore(ctx context.Context, id string, name string) error {
	query := r.table(`
		UPDATE %s
		SET name = ? where id = ?
	`)

	_, err := r.db.Exec(ctx, query, name, id)

	return err
}

func (r *StoreRepository) FindOneByID(ctx context.Context, id string) (*domain.Store, error) {
	query := r.table(`
		SELECT id, name FROM %s
		WHERE id = ?
		LIMIT 1
	`)

	results, err := r.db.Engine(ctx).Query(query, id)
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

func (r *StoreRepository) FindAll(ctx context.Context) ([]*domain.Store, error) {
	panic("unimplemented")
}

func (StoreRepository) table(query string) string {
	return fmt.Sprintf(query, StoreTable)
}
