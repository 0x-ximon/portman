-- +goose Up
-- +goose StatementBegin
CREATE TYPE USER_ROLE AS ENUM ('REGULAR', 'AUTOMATED', 'ADMINISTRATOR');

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),

    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email_address TEXT NOT NULL UNIQUE,
    wallet_address TEXT NOT NULL UNIQUE,

    password TEXT NOT NULL,
    api_key TEXT DEFAULT NULL,
    role USER_ROLE NOT NULL DEFAULT 'REGULAR',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE users;
DROP TYPE USER_ROLE;
-- +goose StatementEnd
