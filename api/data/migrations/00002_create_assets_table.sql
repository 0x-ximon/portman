-- +goose Up
-- +goose StatementBegin
CREATE TYPE ASSET_CLASS AS ENUM ('CRYPTO', 'STOCK', 'FIAT', 'COMMODITY');

CREATE TABLE assets (
      id SERIAL PRIMARY KEY,

      name TEXT NOT NULL,
      symbol TEXT NOT NULL UNIQUE ,
      decimals INTEGER NOT NULL,
      asset_class ASSET_CLASS NOT NULL,

      created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
      updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
      deleted_at TIMESTAMPTZ DEFAULT NULL
  )
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE assets;
DROP TYPE ASSET_CLASS;
-- +goose StatementEnd
