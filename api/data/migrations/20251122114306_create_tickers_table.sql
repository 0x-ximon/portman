-- +goose Up
-- +goose StatementBegin
CREATE TYPE TICKER_STATUS AS ENUM ('OPEN', 'CLOSED', 'SUSPENDED');

CREATE TABLE
  tickers (
    ID SERIAL PRIMARY KEY,

    base TEXT NOT NULL,
    quote TEXT NOT NULL,
    symbol TEXT UNIQUE NOT NULL,

    ask NUMERIC(10, 6) NOT NULL DEFAULT 0.0,
    bid NUMERIC(10, 6) NOT NULL DEFAULT 0.0,
    last NUMERIC(10, 6) NOT NULL DEFAULT 0.0,
    status TICKER_STATUS NOT NULL DEFAULT 'CLOSED'
  )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tickers;
DROP TYPE TICKER_STATUS;
-- +goose StatementEnd
