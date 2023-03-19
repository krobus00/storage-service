-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS object_types (
    id varchar(36) NOT NULL UNIQUE,
    name TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS object_types;
-- +goose StatementEnd
