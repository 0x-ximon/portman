-- name: CreateTicker :one
INSERT INTO tickers (
    symbol, base, quote, status
) VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetTicker :one
SELECT * FROM tickers
WHERE ID = $1 LIMIT 1;

-- name: ListTickers :many
SELECT * FROM tickers
ORDER BY symbol;

-- name: DeleteTicker :exec
DELETE FROM tickers
WHERE ID = $1;