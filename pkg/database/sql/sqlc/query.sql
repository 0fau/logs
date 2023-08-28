-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByToken :one
SELECT *
FROM users
WHERE access_token = $1
LIMIT 1;

-- name: SetAccessToken :exec
UPDATE users
SET access_token = $2
WHERE id = $1;

-- name: UpsertUser :one
INSERT
INTO users (discord_id, discord_name)
VALUES ($1, $2)
ON CONFLICT (discord_id)
DO UPDATE SET discord_name = excluded.discord_name
RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;