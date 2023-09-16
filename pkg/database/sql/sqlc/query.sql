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
INTO encounters (uploaded_by, raid, date, visibility, duration, total_damage_dealt, cleared, local_player)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: InsertEntity :one
INSERT
INTO entities (encounter, class, enttype, name, damage, dps)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetEncounter :one
SELECT *
FROM encounters
WHERE id = $1
LIMIT 1;

-- name: ListRecentEncounters :many
SELECT *
FROM encounters
ORDER BY date DESC
LIMIT 5;

-- name: GetEntities :many
SELECT *
FROM entities
WHERE encounter = $1;

-- name: InsertBuff :copyfrom
INSERT INTO buffs (encounter, player, buff_id, percent, damage)
VALUES ($1, $2, $3, $4, $5);

-- name: InsertSkill :copyfrom
INSERT INTO skills (encounter, player, skill_id, casts, crits, dps, hits, max_damage, total_damage, name)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10);

-- name: GetBuffs :many
SELECT *
FROM buffs
WHERE encounter = $1;

-- name: GetSkills :many
SELECT *
FROM skills
WHERE encounter = $1;