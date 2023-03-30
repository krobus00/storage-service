-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS objects (
    id varchar(36) UNIQUE,
    file_name text NOT NULL,
	key text NOT NULL UNIQUE,
    uploaded_by varchar(36) NOT NULL,
    is_public boolean NOT NULL,
    type_id varchar(36) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_type FOREIGN KEY(type_id) REFERENCES object_types(id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS objects;
-- +goose StatementEnd
