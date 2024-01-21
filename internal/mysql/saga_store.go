package mysql

import (
	"context"
	"fmt"

	"github.com/htquangg/microservices-poc/internal/sec"
	"github.com/htquangg/microservices-poc/pkg/converter"
	"github.com/htquangg/microservices-poc/pkg/database"
)

const SagaTable = "sagas"

var _ sec.SagaStore = (*SagaStore)(nil)

type SagaStore struct {
	db database.DB
}

func NewSagaStore(db database.DB) SagaStore {
	return SagaStore{
		db: db,
	}
}

func (s SagaStore) Load(ctx context.Context, sagaName string, sagaID string) (*sec.SagaContext[[]byte], error) {
	query := s.table("SELECT data, step, done, compensating from %s where id = ? AND name = ?")

	results, err := s.db.Engine(ctx).Query(query, sagaID, sagaName)
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, ErrRecordNotFound
	}

	return &sec.SagaContext[[]byte]{
		ID:           sagaID,
		Data:         results[0]["data"],
		Step:         converter.StringToInt(string(results[0]["step"])),
		Done:         converter.StringToBoolean(string(results[0]["done"])),
		Compensating: converter.StringToBoolean(string(results[0]["compensating"])),
	}, nil
}

func (s SagaStore) Save(ctx context.Context, sagaName string, sagaCtx *sec.SagaContext[[]byte]) error {
	query := s.table(`
		INSERT INTO %s (id, name, data, step, done, compensating)
		VALUES (?, ?, ?, ?, ?, ?) as new
		ON DUPLICATE KEY
		UPDATE
			data=new.data,
			step=new.step,
			done=new.done,
			compensating=new.compensating
	`)

	_, err := s.db.Exec(ctx, query, sagaCtx.ID, sagaName, sagaCtx.Data, sagaCtx.Step, sagaCtx.Done, sagaCtx.Compensating)

	return err
}

func (s SagaStore) table(query string) string {
	return fmt.Sprintf(query, SagaTable)
}
