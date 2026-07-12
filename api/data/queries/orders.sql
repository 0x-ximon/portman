-- name: CreateOrder :one
INSERT INTO orders (user_id, ticker_id, price, quantity, order_side, order_mode)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE ID = $1 and user_id = $2 LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders
WHERE user_id = $1
ORDER BY created_at;

-- name: UpdateOrder :exec
UPDATE orders
SET order_status = $2, updated_at = now()
WHERE ID = $1;
