package mysql

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/es"
	"github.com/htquangg/microservices-poc/internal/registry"
	"github.com/htquangg/microservices-poc/pkg/converter"
	"github.com/htquangg/microservices-poc/pkg/database"
)

const EventTable = "events"

type (
	EventStore struct {
		db       database.DB
		registry registry.Registry
	}

	aggregateEvent struct {
		id         string
		name       string
		payload    ddd.EventPayload
		occurredAt time.Time
		aggregate  es.EventSourcedAggregate
		version    int
	}
)

var (
	_ es.AggregateStore  = (*EventStore)(nil)
	_ ddd.AggregateEvent = (*aggregateEvent)(nil)
)

func NewEventStore(db database.DB, registry registry.Registry) EventStore {
	return EventStore{
		db:       db,
		registry: registry,
	}
}

func (s EventStore) Load(ctx context.Context, aggregate es.EventSourcedAggregate) error {
	query := s.table(
		"SELECT stream_version, event_id, event_name, event_data, occurred_at FROM %s WHERE stream_id = ? AND stream_name = ? AND stream_version > ? ORDER BY stream_version ASC",
	)

	events, err := s.db.Engine(ctx).Query(query, aggregate.ID(), aggregate.AggregateName(), aggregate.Version())
	if err != nil {
		return err
	}

	for _, event := range events {
		eventID := string(event["event_id"])
		eventName := string(event["event_name"])
		payloadData := event["event_data"]
		aggregateVersion := converter.StringToInt(string(event["stream_version"]))
		occurredAt, err := time.ParseInLocation("2006-01-02 15:04:05", string(event["occurred_at"]), time.Local)
		if err != nil {
			return err
		}

		var payload interface{}
		payload, err = s.registry.Deserialize(eventName, payloadData)
		if err != nil {
			return err
		}

		event := aggregateEvent{
			id:         eventID,
			name:       eventName,
			payload:    payload,
			aggregate:  aggregate,
			version:    aggregateVersion,
			occurredAt: occurredAt,
		}

		if err = es.LoadEvent(aggregate, event); err != nil {
			return err
		}

	}

	return nil
}

func (s EventStore) Save(ctx context.Context, aggregate es.EventSourcedAggregate) error {
	buf := &strings.Builder{}
	fmt.Fprint(buf,
		s.table(`
			INSERT INTO %s (stream_id, stream_name, stream_version, event_id, event_name, event_data, occurred_at) VALUES
	`),
	)

	vals := make([]interface{}, 1, len(aggregate.Events())*7+1)

	for idx, event := range aggregate.Events() {
		payloadData, err := s.registry.Serialize(event.EventName(), event.Payload())
		if err != nil {
			return err
		}

		buf.WriteString("(?, ?, ?, ?, ?, ?, ?)")
		vals = append(vals,
			aggregate.ID(),
			aggregate.AggregateName(),
			event.AggregateVersion(),
			event.ID(),
			event.EventName(),
			payloadData,
			event.OccurredAt(),
		)

		// trim the last ","
		if idx != len(aggregate.Events())-1 {
			buf.WriteString(",")
		}
	}

	vals[0] = buf.String()

	_, err := s.db.Exec(ctx, vals...)

	return err
}

func (e aggregateEvent) ID() string {
	return e.id
}

func (e aggregateEvent) EventName() string {
	return e.name
}

func (e aggregateEvent) Payload() ddd.EventPayload {
	return e.payload
}

func (e aggregateEvent) Metadata() ddd.Metadata {
	return ddd.Metadata{}
}

func (e aggregateEvent) OccurredAt() time.Time {
	return e.occurredAt
}

func (e aggregateEvent) AggregateID() string {
	return e.aggregate.ID()
}

func (e aggregateEvent) AggregateName() string {
	return e.aggregate.AggregateName()
}

func (e aggregateEvent) AggregateVersion() int {
	return e.aggregate.Version()
}

func (s EventStore) table(query string) string {
	return fmt.Sprintf(query, EventTable)
}
