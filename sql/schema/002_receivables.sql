-- +goose Up
CREATE TABLE receivables
(
    id              UUID PRIMARY KEY,
    created_at      TIMESTAMP NOT NULL,
    updated_at      TIMESTAMP NOT NULL,
    userid          UUID      NOT NULL REFERENCES users (id),
    date            TIMESTAMP NOT NULL,
    amount_total    FLOAT     NOT NULL,
    amount_received FLOAT     NOT NULL,
    amount_left     FLOAT     NOT NULL,
    status          TEXT      NOT NULL
);

-- +goose Down
DROP TABLE receivables;