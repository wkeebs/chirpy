-- name: CreateChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(), NOW(), NOW(), $1, $2
)
RETURNING *;

-- name: GetChirp :one
SELECT * FROM chirps WHERE chirps.id = $1;

-- name: GetAllChirps :many
SELECT * FROM chirps ORDER BY chirps.created_at ASC;