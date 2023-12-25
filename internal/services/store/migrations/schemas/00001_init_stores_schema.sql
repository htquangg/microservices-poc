-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS stores
(
    id          BIGINT(20) PRIMARY KEY              NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,

    name        VARCHAR(64)                         NOT NULL
);

CREATE TABLE IF NOT EXISTS products
(
    id          BIGINT(20) PRIMARY KEY              NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,

    store_id    BIGINT(20)                          NOT NULL,
    name        VARCHAR(64)                         NOT NULL,
    description VARCHAR(2048)                       NULL,
    sku         VARCHAR(2048)                       NOT NULL,
    price       DECIMAL(19, 4)                      NOT NULL
);

CREATE INDEX idx_store_id ON products (store_id);

CREATE TABLE IF NOT EXISTS outboxes
(
    id           BIGINT(20) PRIMARY KEY NOT NULL,
    subject      VARCHAR(255)           NOT NULL,
    name         VARCHAR(255)           NOT NULL,
    data         BLOB                   NOT NULL,
    metadata     BLOB                   NOT NULL,
    sent_at      DATETIME               NOT NULL,
    published_at DATETIME
);

CREATE TABLE events
(
    stream_id      BIGINT(20)                          NOT NULL,
    stream_name    VARCHAR(512)                        NOT NULL,
    stream_version int                                 NOT NULL,
    event_id       BIGINT(20)                          NOT NULL,
    event_name     VARCHAR(512)                        NOT NULL,
    event_data     BLOB                                NOT NULL,
    occurred_at    DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (stream_id, stream_name, stream_version)
);

CREATE TABLE snapshots
(
    stream_id      BIGINT(20)                          NOT NULL,
    stream_name    VARCHAR(512)                        NOT NULL,
    stream_version int                                 NOT NULL,
    snapshot_name  VARCHAR(512)                        NOT NULL,
    snapshot_data  BLOB                                NOT NULL,
    updated_at     DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    PRIMARY KEY (stream_id, stream_name)
);


-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
DROP TABLE IF EXISTS snapshots;

DROP TABLE IF EXISTS events;

DROP TABLE IF EXISTS outboxes;

DROP INDEX idx_store_id ON products;

DROP TABLE IF EXISTS products;

DROP TABLE IF EXISTS stores;
