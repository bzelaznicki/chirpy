-- name: AddRefreshToken :one
INSERT INTO refresh_tokens (token, created_at, updated_at, user_id, expires_at)
VALUES ($1, NOW(), NOW(), $2, $3)
RETURNING *;

-- name: LookUpRefreshToken :one
SELECT * FROM refresh_tokens where token = $1 AND revoked_at IS NULL AND expires_at > NOW();

-- name: RevokeRefreshToken :exec
UPDATE refresh_tokens
SET revoked_at = NOW(),
updated_at = NOW()
WHERE token = $1;