-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING *;

-- name: GetUserByID :one
SELECT * FROM users WHERE users.id = $1;

-- name: GetAllUsers :many
SELECT * FROM users ORDER BY users.created_at ASC;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users WHERE users.email = $1;

-- name: UpdateUser :one
UPDATE users
SET 
    email = $2,
    hashed_password = $3,
    updated_at = NOW()
WHERE 
    users.id = $1
RETURNING *;

-- name: UpgradeUserToPremium :one
UPDATE users
SET is_premium = true
WHERE users.id = $1
RETURNING *;