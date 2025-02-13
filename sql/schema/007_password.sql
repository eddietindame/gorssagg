-- +goose Up

ALTER TABLE users ADD COLUMN password TEXT NOT NULL DEFAULT 'password' CHECK (length(password) > 0);

-- +goose Down

ALTER TABLE users DROP COLUMN password;
