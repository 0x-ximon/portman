-- +goose Up
-- +goose StatementBegin
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE ROLE AS ENUM ('REGULAR', 'AUTOMATED', 'ADMINISTRATOR');

CREATE TABLE
  users (
    ID UUID PRIMARY KEY DEFAULT uuid_generate_v4(),

    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    phone_number TEXT NOT NULL UNIQUE,
    email_address TEXT NOT NULL UNIQUE,
    wallet_address TEXT NOT NULL UNIQUE,

    free_balance NUMERIC NOT NULL DEFAULT 0.0,
    frozen_balance NUMERIC NOT NULL DEFAULT 0.0,

    password TEXT NOT NULL,
    api_key TEXT DEFAULT NULL,
    role ROLE NOT NULL DEFAULT 'REGULAR',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
  )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TYPE ROLE;
-- +goose StatementEnd
