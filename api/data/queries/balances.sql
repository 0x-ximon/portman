-- name: CreateBalance :one
INSERT INTO balances (user_id, asset_id, available, locked)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetBalance :one
SELECT * FROM balances
WHERE user_id = $1 AND asset_id = $2 LIMIT 1;

-- name: ListBalances :many
SELECT * FROM balances
WHERE user_id = $1
ORDER BY asset_id;

-- name: UpdateBalance :exec
UPDATE balances
SET available = $3, locked = $4
WHERE user_id = $1 AND asset_id = $2;
