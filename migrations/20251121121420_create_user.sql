-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id     TEXT PRIMARY KEY,
    username    TEXT NOT NULL,
    is_active   BOOLEAN DEFAULT TRUE,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd
