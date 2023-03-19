-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS object_whitelist_types (
    id varchar(36) UNIQUE,
    type_id varchar(36) NOT NULL,
    extension varchar(5) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_type FOREIGN KEY(type_id) REFERENCES object_types(id),
    CONSTRAINT uniq_type_id_and_extensions UNIQUE (type_id, extension)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS object_whitelist_types;
-- +goose StatementEnd
