-- name: GetUser :one
SELECT *
FROM users
WHERE discord_tag = $1
LIMIT 1;

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetUserByToken :one
SELECT *
FROM users
WHERE access_token = $1
LIMIT 1;

-- name: SetUsername :exec
UPDATE users
SET username = $2
WHERE id = $1;

-- name: SetAccessToken :exec
UPDATE users
SET access_token = $2
WHERE id = $1;

-- name: UpsertUser :one
INSERT
INTO users (discord_id, discord_tag, avatar, settings)
VALUES ($1, $2, $3, $4)
ON CONFLICT (discord_id)
    DO UPDATE SET discord_tag = excluded.discord_tag,
                  avatar      = excluded.avatar
RETURNING *;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;

-- name: InsertEncounter :one
INSERT
INTO encounters (uploaded_by, settings, tags, header, data, difficulty, boss, date, duration, local_player)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
RETURNING id;

-- name: InsertEntity :one
INSERT
INTO players (encounter, class, enttype, name, damage, dps, dead, fields)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetEncounter :one
SELECT *
FROM encounters
WHERE id = $1
LIMIT 1;

-- name: ListRecentEncounters :many
SELECT id,
       difficulty,
       uploaded_by,
       uploaded_at,
       settings,
       tags,
       header,
       boss,
       date,
       duration,
       local_player
FROM encounters
WHERE (sqlc.narg('date')::TIMESTAMP IS NULL
    OR sqlc.narg('date') >= date)
  AND (sqlc.narg('id') < id
    OR sqlc.narg('id') IS NULL)
  AND (sqlc.narg('user')::UUID IS NULL
    OR sqlc.narg('user') = uploaded_by)
ORDER BY date DESC, id ASC
LIMIT 5;

-- name: GetEntities :many
SELECT *
FROM players
WHERE encounter = $1;

-- name: InsertSkill :copyfrom
INSERT INTO skills (encounter, player, skill_id, dps, damage, name, tripods, fields)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetData :one
SELECT data
FROM encounters
WHERE id = $1;

-- name: GetSkills :many
SELECT *
FROM skills
WHERE encounter = $1;