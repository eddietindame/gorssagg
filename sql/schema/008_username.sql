-- +goose Up

ALTER TABLE users DROP COLUMN name;
ALTER TABLE users ADD COLUMN username TEXT NOT NULL UNIQUE CHECK (length(username) > 0) DEFAULT (gen_random_uuid());

-- +goose Down

ALTER TABLE users ADD COLUMN name TEXT NOT NULL DEFAULT (gen_random_uuid());
ALTER TABLE users DROP COLUMN username;
