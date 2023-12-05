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
	fmt.Fprintf(buf, "INSERT INTO %s (id, subject, name, data, metadata, sent_at) VALUES ", OUTBOX_TABLE)

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

func (s OutboxStore) FindUnpublished(ctx context.Context, currentOffset string, limit int) ([]am.Message, error) {
	query := fmt.Sprintf("SELECT id, subject, name, data, metadata, sent_at FROM %s "+
		"WHERE id > ? AND published_at IS NULL ORDER BY id ASC LIMIT ?", OUTBOX_TABLE)

	results, err := s.db.Engine(ctx).Query(query, currentOffset, limit)
	if err != nil {
		return nil, err
	}

	msgs := make([]am.Message, 0, len(results))

	for _, result := range results {
		tp, err := time.ParseInLocation("2006-01-02 15:04:05", string(result["sent_at"]), time.Local)
		if err != nil {
			return nil, err
		}

		msg := outboxMessage{
			id:      string(result["id"]),
			subject: string(result["subject"]),
			name:    string(result["name"]),
			data:    result["data"],
			sentAt:  tp,
		}

		if err = json.Unmarshal(result["metadata"], &msg.metadata); err != nil {
			return nil, err
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

func (s OutboxStore) MarkPublished(ctx context.Context, currentOffset string, lastOffset string) error {
	query := fmt.Sprintf("UPDATE %s SET published_at = CURRENT_TIMESTAMP where id >= ? and id <= ?", OUTBOX_TABLE)

	_, err := s.db.Exec(ctx, query, currentOffset, lastOffset)

	return err
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
