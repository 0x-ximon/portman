-- +goose Up
-- +goose StatementBegin
CREATE TYPE TICKER_STATUS AS ENUM ('OPEN', 'CLOSED', 'SUSPENDED');

CREATE TABLE
  tickers (
    ID SERIAL PRIMARY KEY,

    base TEXT NOT NULL,
    quote TEXT NOT NULL,
    symbol TEXT UNIQUE NOT NULL,
    status TICKER_STATUS NOT NULL DEFAULT 'CLOSED'
  )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tickers;
DROP TYPE TICKER_STATUS;
-- +goose StatementEnd
