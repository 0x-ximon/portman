-- +goose Up
-- +goose StatementBegin
CREATE FUNCTION to_precision(decimals INTEGER) RETURNS NUMERIC AS $$
BEGIN
    RETURN 10::NUMERIC ^ (-decimals);
END;
$$ LANGUAGE plpgsql;

CREATE TYPE TICKER_STATUS AS ENUM ('OPEN', 'CLOSED', 'SUSPENDED');

CREATE TABLE tickers (
    id BIGSERIAL PRIMARY KEY,
    UNIQUE (base_asset, quote_asset),

    symbol TEXT UNIQUE NOT NULL,
    lot_size NUMERIC NOT NULL,
    tick_size NUMERIC NOT NULL,
    ask NUMERIC NOT NULL DEFAULT 0,
    bid NUMERIC NOT NULL DEFAULT 0,

    base_asset INTEGER REFERENCES assets(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
    quote_asset INTEGER REFERENCES assets(id) ON DELETE CASCADE ON UPDATE CASCADE NOT NULL,
    status TICKER_STATUS NOT NULL DEFAULT 'CLOSED',

    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    deleted_at TIMESTAMPTZ DEFAULT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tickers;
DROP TYPE TICKER_STATUS;
DROP FUNCTION to_precision;
-- +goose StatementEnd
