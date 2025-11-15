-- name: GetUser :one
SELECT * FROM users
WHERE ID = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users;

-- name: CreateUser :one
INSERT INTO users (
    first_name, last_name, email_address, wallet_address, role, password
) VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: DeleteUser :exec
DELETE FROM users
WHERE ID = $1;