-- name: CreateTicker :one
INSERT INTO tickers (symbol, lot_size, tick_size, base_asset, quote_asset)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetTicker :one
SELECT * FROM tickers
WHERE ID = $1 LIMIT 1;

-- name: FindTickerBySymbol :one
SELECT * FROM tickers
WHERE symbol = $1 LIMIT 1;

-- name: UpdateTicker :one
UPDATE tickers
SET lot_size = $2, tick_size = $3
WHERE ID = $1
RETURNING *;

-- name: UpdateTickerQuotes :one
UPDATE tickers
SET ask = $2, bid = $3
WHERE ID = $1
RETURNING *;

-- name: UpdateTickerStatus :one
UPDATE tickers
SET status = $2
WHERE ID = $1
RETURNING *;

-- name: ListTickers :many
SELECT * FROM tickers
ORDER BY symbol;

-- name: DeleteTicker :exec
DELETE FROM tickers
WHERE ID = $1;
