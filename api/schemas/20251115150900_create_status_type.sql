-- +goose Up
-- +goose StatementBegin
CREATE TYPE STATUS AS ENUM ('OPEN', 'CLOSED', 'SUSPENDED');

CREATE TABLE
  tickers (
    ID SERIAL PRIMARY KEY,
    symbol TEXT UNIQUE NOT NULL,
    base TEXT NOT NULL,
    quote TEXT NOT NULL,
    status STATUS NOT NULL
  )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE tickers;
DROP TYPE STATUS;
-- +goose StatementEnd