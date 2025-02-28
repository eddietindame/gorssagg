-- +goose Up

CREATE EXTENSION IF NOT EXISTS citext;

CREATE DOMAIN username AS CITEXT
CHECK (
  VALUE ~* '^[A-Za-z][A-Za-z0-9_]{7,29}$'
);

CREATE DOMAIN email_address AS CITEXT
CHECK (
  VALUE ~* '^[\w.%+-]+@[a-zA-Z0-9]+(\.[a-zA-Z0-9-]+)*\.[a-zA-Z]{2,}$'
);

CREATE TABLE users(
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  email email_address NOT NULL UNIQUE,
  username username NOT NULL UNIQUE,
  password TEXT NOT NULL CHECK (length(password) > 0)
);

-- +goose Down

DROP TABLE users;
DROP DOMAIN IF EXISTS username;
DROP DOMAIN IF EXISTS email_address;
DROP EXTENSION IF EXISTS citext;
