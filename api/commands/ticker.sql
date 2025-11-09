-- name: GetTicker :one
SELECT * FROM tickers
WHERE ID = $1 LIMIT 1;

-- name: ListTickers :many
SELECT * FROM tickers
ORDER BY name;

-- name: CreateTicker :one
INSERT INTO tickers (
    symbol, name
) VALUES (
  $1, $2
)
RETURNING *;

-- name: DeleteTicker :exec
DELETE FROM tickers
WHERE ID = $1;