-- name: CreateOrder :one
INSERT INTO orders (user_id, ticker_symbol, price, quantity, side, type)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE ID = $1 and user_id = $2
LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders
WHERE user_id = $1;
