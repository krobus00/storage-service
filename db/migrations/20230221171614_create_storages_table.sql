-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS storages (
    id varchar(36) UNIQUE,
    file_name text NOT NULL,
	object_key text NOT NULL UNIQUE,
    is_public boolean NOT NULL,
    uploaded_by varchar(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS storages;
-- +goose StatementEnd
