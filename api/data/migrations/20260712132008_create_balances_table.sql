-- +goose Up
-- +goose StatementBegin
CREATE TABLE balances (
    PRIMARY KEY (user_id, asset_id),
    user_id UUID REFERENCES users(id) NOT NULL,
    asset_id INTEGER REFERENCES assets(id) NOT NULL,

    available NUMERIC NOT NULL DEFAULT 0,
    locked NUMERIC NOT NULL DEFAULT 0
);

ALTER TABLE balances
ADD CONSTRAINT validate_available CHECK (available >= 0),
ADD CONSTRAINT validate_locked CHECK (locked >= 0);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE balances;
-- +goose StatementEnd
