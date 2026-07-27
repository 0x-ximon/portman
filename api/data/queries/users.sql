-- name: CreateUser :one
INSERT INTO users (first_name, last_name, email_address, wallet_address, password, api_key)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users
WHERE ID = $1 LIMIT 1;

-- name: FindUserByEmail :one
SELECT * FROM users
WHERE email_address = $1 LIMIT 1;

-- name: FindUserByApiKey :one
SELECT * FROM users
WHERE api_key = $1 LIMIT 1;

-- name: ListUsers :many
SELECT * FROM users;

-- name: UpdateUser :exec
UPDATE users
SET first_name = $2, last_name = $3, email_address = $4, wallet_address = $5, password = $6
WHERE ID = $1;

-- name: UpdateUserRole :exec
UPDATE users
SET role = $2
WHERE ID = $1;

-- name: DeleteUser :exec
UPDATE users
SET deleted_at = now()
WHERE ID = $1;
