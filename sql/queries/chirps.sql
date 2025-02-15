-- name: PostChirp :one
INSERT INTO chirps (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
) RETURNING *;

-- name: GetChirps :many
SELECT * FROM chirps ORDER BY created_at ASC;

-- name: GetChirpsByUser :many
SELECT * FROM chirps where user_id = $1 ORDER BY created_at ASC;

-- name: GetSingleChirpByUUID :one

SELECT * FROM chirps WHERE id = $1;

-- name: DeleteChirp :exec

DELETE FROM chirps WHERE id = $1;