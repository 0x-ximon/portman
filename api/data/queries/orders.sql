-- name: CreateOrder :one
INSERT INTO orders (buyer_id, seller_id, price, quantity, side, type, status)
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING *;

-- name: GetOrder :one
SELECT * FROM orders
WHERE ID = $1 and buyer_id = $2 or seller_id = $2
LIMIT 1;

-- name: ListOrders :many
SELECT * FROM orders
WHERE buyer_id = $1 or seller_id = $1;
