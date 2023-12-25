package domain

import (
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/es"
	"github.com/stackus/errors"
)

const StoreAggregate = "stores.Store"

var ErrStoreNameIsBlank = errors.Wrap(errors.ErrBadRequest, "the store name cannot be blank")

var _ interface {
	es.EventApplier
	es.Snapshotter
} = (*StoreES)(nil)

type StoreES struct {
	es.Aggregate
	name string
}

// Key implements registry.Registerable
func (StoreES) Key() string {
	return StoreAggregate
}

func NewStore(id string) *StoreES {
	return &StoreES{
		Aggregate: es.NewAggregate(id, StoreAggregate),
	}
}

func (s StoreES) Name() string {
	return s.name
}

func (s *StoreES) Init(name string) (ddd.Event, error) {
	if name == "" {
		return nil, ErrStoreNameIsBlank
	}

	s.AddEvent(StoreCreatedEvent, &StoreCreated{
		Name: name,
	})

	return ddd.NewEvent(StoreCreatedEvent, s), nil
}

func (s *StoreES) Rebrand(name string) (ddd.Event, error) {
	if name == "" {
		return nil, ErrStoreNameIsBlank
	}

	s.AddEvent(StoreRebrandedEvent, &StoreRebranded{
		Name: name,
	})

	return ddd.NewEvent(StoreRebrandedEvent, s), nil
}

// ApplyEvent implements es.EventApplier
func (s *StoreES) ApplyEvent(event ddd.Event) error {
	switch payload := event.Payload().(type) {
	case *StoreCreated:
		s.name = payload.Name
	case *StoreRebranded:
		s.name = payload.Name
	default:
		return errors.ErrInternal.Msgf("%T received the event %s with unexpected payload %T", s, event.EventName(), payload)
	}

	return nil
}

// ApplySnapshot implements es.Snapshotter
func (s *StoreES) ApplySnapshot(snapshot es.Snapshot) error {
	switch ss := snapshot.(type) {
	case *StoreV1:
		s.name = ss.Name
	default:
		return errors.ErrInternal.Msgf("%T received the unexpected snapshot %T", s, snapshot)

	}

	return nil
}

// ToSnapshot implements es.Snapshotter
func (s *StoreES) ToSnapshot() es.Snapshot {
	return StoreV1{
		Name: s.name,
	}
}
