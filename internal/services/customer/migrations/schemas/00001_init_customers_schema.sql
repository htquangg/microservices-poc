-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd

CREATE TABLE IF NOT EXISTS customers
(
    id         BIGINT(20) PRIMARY KEY              NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,

    name       VARCHAR(64)                         NOT NULL,
    email      VARCHAR(64)                         NULL,
    phone      VARCHAR(16)                         NOT NULL
);

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

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd

DROP TABLE IF EXISTS outboxes;

DROP TABLE IF EXISTS customers;
