-- name: GetUser :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: CreateUser :one
INSERT INTO users (discord_id, discord_name)
VALUES ($1, $2)
RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;