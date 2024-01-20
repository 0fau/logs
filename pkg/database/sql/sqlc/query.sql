-- name: GetUser :one
SELECT *
FROM users
WHERE discord_tag = $1
LIMIT 1;

-- name: UpdateAvatar :exec
UPDATE users
SET avatar = $2
WHERE id = $1;

-- ROLES

-- name: GetUserByID :one
SELECT *
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetRoles :one
SELECT roles
FROM users
WHERE id = $1
LIMIT 1;

-- name: GetRolesByDiscordID :one
SELECT id, roles
FROM users
WHERE discord_id = $1
LIMIT 1;

-- name: UpdateRoles :exec
UPDATE users
SET roles = $2
WHERE id = $1;

-- name: Whitelist :exec
INSERT INTO whitelist (discord, role)
VALUES ($1, $2)
ON CONFLICT DO NOTHING;

-- name: FetchWhitelist :one
DELETE
FROM whitelist
WHERE discord = $1
RETURNING role;

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
INTO users (discord_id, discord_tag, roles, avatar, settings)
VALUES ($1, $2, $3, $4, $5)
ON CONFLICT (discord_id)
    DO UPDATE SET discord_tag = excluded.discord_tag,
                  roles       = excluded.roles
RETURNING *;

-- name: SetUserRoles :exec
UPDATE users
SET roles = $2
WHERE discord_tag = $1;

-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1;

-- name: InsertEncounter :one
INSERT
INTO encounters (uploaded_by, settings, tags, header, data, difficulty, boss, date, duration, local_player, unique_hash,
                 unique_group)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id;

-- name: GetUniqueGroup :one
SELECT unique_group
FROM encounters
WHERE unique_hash = $1
  AND unique_group = id
  AND (date + interval '120 seconds') >= $2
  AND (date - interval '120 seconds') <= $2;

-- name: UpdateUniqueGroup :exec
UPDATE encounters
SET unique_group = $2
WHERE id = $1;

-- name: UpsertEncounterGroup :exec
INSERT
INTO grouped_encounters (group_id, uploaders)
VALUES ($1, ARRAY [$2::UUID])
ON CONFLICT (group_id)
    DO UPDATE SET uploaders = array_append(grouped_encounters.uploaders, $2);

-- name: InsertPlayer :copyfrom
INSERT
INTO players (encounter, class, name, dead, data, dps, place)
VALUES ($1, $2, $3, $4, $5, $6, $7);

-- name: InsertPlayerInternal :exec
INSERT
INTO players (encounter, class, name, dead, data, dps, place)
VALUES ($1, $2, $3, $4, $5, $6, $7)
ON CONFLICT (encounter, name)
    DO UPDATE SET data = excluded.data;

-- name: ProcessEncounter :one
UPDATE encounters
SET header      = $2,
    data        = $3,
    unique_hash = $4
WHERE id = $1
RETURNING uploaded_by;

-- name: GetEncounter :one
SELECT e.*,
       u.discord_id,
       u.discord_tag,
       u.username,
       u.avatar
FROM encounters e
         JOIN users u ON e.uploaded_by = u.id
WHERE e.id = $1
LIMIT 1;

-- name: DeleteEncounter :exec
DELETE
FROM encounters
WHERE id = $1;

-- name: ListEncounters :many
SELECT id
FROM encounters
ORDER BY uploaded_at;

-- name: GetHeader :one
SELECT boss, difficulty, uploaded_by, date, duration, header
FROM encounters
WHERE id = $1
LIMIT 1;