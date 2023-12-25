package domain

const (
	StoreCreatedEvent   = "stores.StoreCreated"
	StoreRebrandedEvent = "stores.StoreRebranded"
)

type StoreCreated struct {
	Name string
}

// Key implements registry.Registrable
func (StoreCreated) Key() string {
	return StoreCreatedEvent
}

type StoreRebranded struct {
	Name string
}

// Key implements registry.Registrable
func (StoreRebranded) Key() string {
	return StoreRebrandedEvent
}
