package mysql

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/htquangg/microservices-poc/internal/am"
	"github.com/htquangg/microservices-poc/internal/tm"
	"github.com/htquangg/microservices-poc/pkg/database"

	mysql_driver "github.com/go-sql-driver/mysql"
)

const INBOX_TABLE = "inboxes"

type InboxStore struct {
	db *database.DB
}

var _ tm.InboxStore = (*InboxStore)(nil)

func NewInboxStore(db *database.DB) InboxStore {
	return InboxStore{
		db: db,
	}
}

func (s InboxStore) Save(ctx context.Context, msg am.IncomingMessage) error {
	query := "INSERT INTO %s (id, subject, name, data, metadata, sent_at, received_at) VALUES (?, ?, ?, ?, ?, ?, ?)"

	metadata, err := json.Marshal(msg.Metadata())
	if err != nil {
		return err
	}

	_, err = s.db.Exec(
		ctx,
		s.table(query),
		msg.ID(),
		msg.Subject(),
		msg.MessageName(),
		msg.Data(),
		metadata,
		msg.SentAt(),
		msg.ReceivedAt(),
	)

	if err != nil {
		var mysqlErr *mysql_driver.MySQLError
		if errors.As(err, &mysqlErr) {
			if mysqlErr.Number == 1062 { // ER_DUP_ENTRY (https://dev.mysql.com/doc/mysql-errors/8.0/en/server-error-reference.html#error_er_dup_entry)
				return tm.ErrDuplicateMessage(msg.ID())
			}
		}
	}

	return err
}

func (s InboxStore) table(query string) string {
	return fmt.Sprintf(query, INBOX_TABLE)
}
