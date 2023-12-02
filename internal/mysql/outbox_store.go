package mysql

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/ddd"
	"github.com/htquangg/microservices-poc/internal/tm"
	"github.com/htquangg/microservices-poc/pkg/database"
)

const OUTBOX_TABLE = "outboxes"

type OutboxStore struct {
	db *database.DB
}

type outboxMessage struct {
	id       string
	subject  string
	name     string
	data     []byte
	metadata ddd.Metadata
	sentAt   time.Time
}

func (outboxMessage) TableName() string {
	return OUTBOX_TABLE
}

var (
	_ tm.OutboxStore = (*OutboxStore)(nil)
	_ am.Message     = (*outboxMessage)(nil)
)

func NewOutboxStore(db *database.DB) *OutboxStore {
	return &OutboxStore{
		db: db,
	}
}

func (s OutboxStore) Save(ctx context.Context, msgs ...am.Message) error {
	if len(msgs) == 0 {
		return nil
	}

	buf := &strings.Builder{}
	buf.WriteString(fmt.Sprintf("INSERT INTO %s (id, subject, name, data, metadata, sent_at) VALUES ", OUTBOX_TABLE))

	vals := make([]interface{}, 1, len(msgs)+1)

	for idx, msg := range msgs {
		metadata, err := json.Marshal(msg.Metadata())
		if err != nil {
			return err
		}

		buf.WriteString("(?, ?, ?, ?, ?, ?)")
		vals = append(vals,
			msg.ID(),
			msg.Subject(),
			msg.MessageName(),
			msg.Data(),
			metadata,
			msg.SentAt(),
		)

		// trim the last ","
		if idx != len(msgs)-1 {
			buf.WriteString(",")
		}
	}

	vals[0] = buf.String()

	_, err := s.db.Exec(ctx, vals...)

	return err
}

func (s OutboxStore) FindUnpublished(ctx context.Context, limit int) ([]am.Message, error) {
	panic("not implemented") // TODO: Implement
}

func (s OutboxStore) MarkPublished(ctx context.Context, ids string) error {
	panic("not implemented") // TODO: Implement
}

func (m outboxMessage) ID() string {
	return m.id
}

func (m outboxMessage) Subject() string {
	return m.subject
}

func (m outboxMessage) MessageName() string {
	return m.name
}

func (m outboxMessage) Metadata() ddd.Metadata {
	return m.metadata
}

func (m outboxMessage) SentAt() time.Time {
	return m.sentAt
}

func (m outboxMessage) Data() []byte {
	return m.data
}
