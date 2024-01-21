-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS sagas
(
    id           BIGINT(20) PRIMARY KEY             NOT NULL,
    updated_at   DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    name         VARCHAR(255)                       NOT NULL,
    data         BLOB                               NOT NULL,
    step         TINYINT                            NOT NULL,
    done         TINYINT                            NOT NULL,
    compensating TINYINT                            NOT NULL
);

CREATE TABLE IF NOT EXISTS inboxes
(
  id           BIGINT(20) PRIMARY KEY NOT NULL,
  subject      VARCHAR(255)           NOT NULL,
  name         VARCHAR(255)           NOT NULL,
  data         BLOB                   NOT NULL,
  metadata     BLOB                   NOT NULL,
  sent_at      DATETIME               NOT NULL,
  received_at  DATETIME
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

DROP TABLE IF EXISTS inboxes;

DROP TABLE IF EXISTS sagas;
