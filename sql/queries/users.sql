-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1
)
RETURNING *;

-- name: GetUser :one
SELECT * FROM users WHERE users.id = $1;

-- name: GetAllUsers :many
SELECT * FROM users ORDER BY users.created_at ASC;

-- name: DeleteAllUsers :exec
DELETE FROM users;