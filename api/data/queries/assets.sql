-- name: CreateAsset :one
INSERT INTO assets (name, symbol, decimals, asset_class)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetAsset :one
SELECT * FROM assets
WHERE id = $1;

-- name: FindAssetBySymbol :one
SELECT * FROM assets
WHERE symbol = $1;

-- name: ListAssets :many
SELECT * FROM assets;

-- name: DeleteAsset :exec
UPDATE assets
SET deleted_at = now()
WHERE id = $1;
