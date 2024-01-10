-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
-- +goose StatementEnd
CREATE TABLE IF NOT EXISTS orders
(
    id          BIGINT(20) PRIMARY KEY             NOT NULL,
    created_at  DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at  DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,

    customer_id BIGINT(20)                         NOT NULL,
    payment_id  BIGINT(20)                         NOT NULL,
    invoice_id  BIGINT(20)                         NOT NULL,
    shopping_id BIGINT(20)                         NOT NULL,
    items       BLOB                               NOT NULL,
    status      VARCHAR(64)                        NOT NULL
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

DROP TABLE IF EXISTS inboxes;

DROP TABLE IF EXISTS orders;
