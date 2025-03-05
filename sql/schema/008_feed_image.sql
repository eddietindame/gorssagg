-- +goose Up

ALTER TABLE feeds ADD COLUMN image TEXT;

-- +goose Down

ALTER TABLE feeds DROP COLUMN image;
