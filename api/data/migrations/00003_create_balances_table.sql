-- +goose Up
-- +goose StatementBegin
CREATE TABLE balances (
    PRIMARY KEY (user_id, asset_id),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
    asset_id INTEGER REFERENCES assets(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,

    available NUMERIC NOT NULL DEFAULT 0,
    locked NUMERIC NOT NULL DEFAULT 0,

    CONSTRAINT validate_available CHECK (available >= 0),
    CONSTRAINT validate_locked CHECK (locked >= 0)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE balances;
-- +goose StatementEnd
