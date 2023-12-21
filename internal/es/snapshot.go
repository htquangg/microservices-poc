package es

import (
	"fmt"
)

type Snapshot interface {
	SnapshotName() string
}

type SnapshotApplier interface {
	ApplySnapshot(snapshot Snapshot) error
}

type Snapshotter interface {
	SnapshotApplier
	ToSnapshot() Snapshot
}

type snapshotLoader interface {
	SnapshotApplier
	VersionSetter
}

func LoadSnapshot(v interface{}, snapshot Snapshot, version int) error {
	agg, ok := v.(snapshotLoader)
	if !ok {
		return fmt.Errorf("%T does not have the methods implement to load snapshots", v)
	}

	if err := agg.ApplySnapshot(snapshot); err != nil {
		return err
	}
	agg.SetVersion(version)

	return nil
}
