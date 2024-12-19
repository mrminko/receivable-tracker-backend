-- name: CreateUser :one
INSERT INTO users (id, name, created_at, updated_at, phone)
VALUES ($1, $2, $3, $4, $5)
    RETURNING *;

-- name: DeleteUser :one
DELETE FROM users WHERE id=$1
    RETURNING *;

-- name: UpdateUser :one
UPDATE users
SET
    name=$2,
    updated_at = NOW(),
    phone = $3
WHERE id=$1
RETURNING id;

-- name: GetUserByID :one
SELECT * FROM users WHERE id = $1;

-- name: GetAllUsers :many
SELECT * FROM users;