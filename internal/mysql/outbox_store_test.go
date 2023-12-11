package mysql

import (
	"context"
	"fmt"
	"testing"

	"github.com/htquangg/microservices-poc/pkg/database"
	"github.com/htquangg/microservices-poc/pkg/utils/testhelper"
	"github.com/stretchr/testify/assert"
)

var dbConfig = &database.Config{
	Port:            3306,
	Host:            "127.0.0.1",
	User:            "root",
	Password:        "toor",
	Schema:          "dev-local-customer-001",
	Charset:         "utf8mb4",
	AutoMigration:   false,
	LogSQL:          false,
	SslMode:         false,
	MaxIdleConns:    1000,
	MaxOpenConns:    100,
	ConnMaxLifetime: 300,
}

func TestOutboxStore_Save(t *testing.T) {
	db, err := testhelper.LoadTestDB(dbConfig)
	if err != nil {
		fmt.Printf("failed to open database connection: %s\n", err)
		return
	}

	ctx := context.Background()
	outboxStore := NewOutboxStore(db)

	err = outboxStore.Save(ctx)
	assert.NoError(t, err)

	msgs, err := outboxStore.FindUnpublished(ctx, "0", 1)
	assert.NoError(t, err)
	assert.Len(t, msgs, 1)
}
