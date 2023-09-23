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

-- name: InsertEncounter :one
INSERT
INTO encounters (uploaded_by, raid, date, visibility, duration, damage, cleared, local_player, fields)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING *;

-- name: InsertEntity :one
INSERT
INTO entities (encounter, class, enttype, name, damage, dps, dead, fields)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetEncounter :one
SELECT *
FROM encounters
WHERE id = $1
LIMIT 1;

-- name: ListRecentEncounters :many
SELECT id,
       uploaded_by,
       raid,
       date,
       visibility,
       duration,
       damage,
       cleared,
       local_player
FROM encounters
ORDER BY date DESC
LIMIT 5;

-- name: GetEntities :many
SELECT *
FROM entities
WHERE encounter = $1;

-- name: InsertSkill :copyfrom
INSERT INTO skills (encounter, player, skill_id, dps, damage, name, tripods, fields)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8);

-- name: GetFields :one
SELECT fields
FROM encounters
WHERE id = $1;

-- name: GetSkills :many
SELECT *
FROM skills
WHERE encounter = $1;