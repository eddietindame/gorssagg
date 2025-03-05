-- +goose Up

ALTER TABLE feeds ADD COLUMN description TEXT;
ALTER TABLE feeds ADD COLUMN language TEXT;

-- +goose Down

ALTER TABLE feeds DROP COLUMN description;
ALTER TABLE feeds DROP COLUMN language;
