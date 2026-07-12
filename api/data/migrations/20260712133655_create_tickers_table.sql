-- +goose Up
-- +goose StatementBegin
CREATE TYPE TICKER_STATUS AS ENUM ('OPEN', 'CLOSED', 'SUSPENDED');

CREATE TABLE tickers (
    id SERIAL PRIMARY KEY,

    symbol TEXT NOT NULL,
    lot_size NUMERIC NOT NULL,
    tick_size NUMERIC NOT NULL,
    ask INTEGER NOT NULL DEFAULT 0,
    bid INTEGER NOT NULL DEFAULT 0,

    base_asset_id TEXT REFERENCES assets(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
    quote_asset_id TEXT REFERENCES assets(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
    ticker_status TICKER_STATUS NOT NULL DEFAULT 'CLOSED',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
)
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tickers;
DROP TYPE TICKER_STATUS;
-- +goose StatementEnd
