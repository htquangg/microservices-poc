package mysql

import (
	"context"
	"fmt"

	"github.com/htquangg/microservices-poc/internal/es"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/pkg/converter"
	"github.com/htquangg/microservices-poc/pkg/database"
)

const (
	SnapshotTable = "snapshots"
	MaxChanges    = 3
)

type (
	SnapshotStore struct {
		es.AggregateStore
		db       database.DB
		registry registry.Registry
	}
)

var _ es.AggregateStore = (*SnapshotStore)(nil)

func NewSnapshotStore(db database.DB, registry registry.Registry) es.AggregateStoreMiddleware {
	snapshots := SnapshotStore{
		db:       db,
		registry: registry,
	}

	return func(store es.AggregateStore) es.AggregateStore {
		snapshots.AggregateStore = store
		return snapshots
	}
}

func (s SnapshotStore) Load(ctx context.Context, aggregate es.EventSourcedAggregate) error {
	query := s.table(
		`
		SELECT stream_version, snapshot_name, snapshot_data
		FROM %s
		WHERE stream_id = ? AND stream_name = ?
		LIMIT 1
	`,
	)

	snapshots, err := s.db.Engine(ctx).Query(query, aggregate.ID(), aggregate.AggregateName())
	if err != nil {
		return err
	}
	if len(snapshots) == 0 {
		return s.AggregateStore.Load(ctx, aggregate)
	}

	v, err := s.registry.Deserialize(
		string(snapshots[0]["snapshot_name"]),
		snapshots[0]["snapshot_data"],
		registry.ValidateImplements((*es.Snapshot)(nil)),
	)
	if err != nil {
		return err
	}

	if err := es.LoadSnapshot(aggregate, v.(es.Snapshot), converter.StringToInt(string(snapshots[0]["snapshot_version"]))); err != nil {
		return err
	}

	return s.AggregateStore.Load(ctx, aggregate)
}

func (s SnapshotStore) Save(ctx context.Context, aggregate es.EventSourcedAggregate) error {
	if err := s.AggregateStore.Save(ctx, aggregate); err != nil {
		return err
	}

	if !s.shouldSnapshot(aggregate) {
		return nil
	}

	query := s.table(
		`
		INSERT INTO %s (stream_id, stream_name, stream_version, snapshot_name, snapshot_data)
		VALUES (?, ?, ?, ?, ?) AS new
		ON DUPLICATE KEY
		UPDATE
			stream_id=new.stream_id,
			stream_name=new.stream_name,
			stream_version=new.stream_version,
			snapshot_name=new.snapshot_name,
			snapshot_data=new.snapshot_data
		`,
	)

	sser, ok := aggregate.(es.Snapshotter)
	if !ok {
		return fmt.Errorf("%T does not implement es.Snapshotter", aggregate)
	}

	snapshot := sser.ToSnapshot()

	data, err := s.registry.Serialize(snapshot.SnapshotName(), snapshot)
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		ctx,
		query,
		aggregate.ID(),
		aggregate.AggregateName(),
		aggregate.PendingVersion(),
		snapshot.SnapshotName(),
		data,
	)

	return err
}

func (s SnapshotStore) shouldSnapshot(aggregate es.EventSourcedAggregate) bool {
	pendingVersion := aggregate.PendingVersion()
	pendingChanges := len(aggregate.Events())

	return pendingVersion >= MaxChanges && ((pendingChanges >= MaxChanges) ||
		(pendingVersion%MaxChanges < pendingChanges) ||
		(pendingVersion%MaxChanges == 0))
}

func (s SnapshotStore) table(query string) string {
	return fmt.Sprintf(query, SnapshotTable)
}
