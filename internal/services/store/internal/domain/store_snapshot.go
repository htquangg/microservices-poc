package domain

type StoreV1 struct {
	Name string
}

// SnapshotName implements es.Snapshot
func (StoreV1) SnapshotName() string {
	return "stores.StoreV1"
}
