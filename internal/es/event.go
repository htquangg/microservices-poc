package es

import (
	"fmt"

	"github.com/htquangg/microservices-poc/internal/ddd"
)

type EventApplier interface {
	ApplyEvent(ddd.Event) error
}

type EventCommitter interface {
	CommitEvents()
}

type eventLoader interface {
	EventApplier
	VersionSetter
}

func LoadEvent(v interface{}, event ddd.AggregateEvent) error {
	agg, ok := v.(eventLoader)
	if !ok {
		return fmt.Errorf("%T does not have the methods implemented to load events", v)
	}

	if err := agg.ApplyEvent(event); err != nil {
		return err
	}
	agg.SetVersion(event.AggregateVersion())

	return nil
}
