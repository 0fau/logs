// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0
// source: query.sql

package sql

import (
	"context"

	structs "github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/jackc/pgx/v5/pgtype"
)

const countClasses = `-- name: CountClasses :many
SELECT (value ->> 'class')::STRING AS class, COUNT(*)
FROM encounters,
     jsonb_each(header -> 'players') AS player
GROUP BY (value ->> 'class')::STRING
`

type CountClassesRow struct {
	Class string
	Count int64
}

func (q *Queries) CountClasses(ctx context.Context) ([]CountClassesRow, error) {
	rows, err := q.db.Query(ctx, countClasses)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []CountClassesRow
	for rows.Next() {
		var i CountClassesRow
		if err := rows.Scan(&i.Class, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const deleteEncounter = `-- name: DeleteEncounter :exec
DELETE
FROM encounters
WHERE id = $1
`

func (q *Queries) DeleteEncounter(ctx context.Context, id int32) error {
	_, err := q.db.Exec(ctx, deleteEncounter, id)
	return err
}

const deleteUser = `-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getData = `-- name: GetData :one
SELECT data
FROM encounters
WHERE id = $1
`

func (q *Queries) GetData(ctx context.Context, id int32) (structs.EncounterData, error) {
	row := q.db.QueryRow(ctx, getData, id)
	var data structs.EncounterData
	err := row.Scan(&data)
	return data, err
}

const getEncounter = `-- name: GetEncounter :one
SELECT e.id, e.uploaded_by, e.uploaded_at, e.settings, e.tags, e.header, e.data, e.unique_hash, e.unique_group, e.boss, e.difficulty, e.date, e.duration, e.local_player,
       u.discord_id,
       u.discord_tag,
       u.username,
       u.avatar
FROM encounters e
         JOIN users u ON e.uploaded_by = u.id
WHERE e.id = $1
LIMIT 1
`

type GetEncounterRow struct {
	ID          int32
	UploadedBy  pgtype.UUID
	UploadedAt  pgtype.Timestamp
	Settings    structs.EncounterSettings
	Tags        []string
	Header      structs.EncounterHeader
	Data        structs.EncounterData
	UniqueHash  string
	UniqueGroup int32
	Boss        string
	Difficulty  string
	Date        pgtype.Timestamp
	Duration    int32
	LocalPlayer string
	DiscordID   string
	DiscordTag  string
	Username    pgtype.Text
	Avatar      pgtype.Text
}

func (q *Queries) GetEncounter(ctx context.Context, id int32) (GetEncounterRow, error) {
	row := q.db.QueryRow(ctx, getEncounter, id)
	var i GetEncounterRow
	err := row.Scan(
		&i.ID,
		&i.UploadedBy,
		&i.UploadedAt,
		&i.Settings,
		&i.Tags,
		&i.Header,
		&i.Data,
		&i.UniqueHash,
		&i.UniqueGroup,
		&i.Boss,
		&i.Difficulty,
		&i.Date,
		&i.Duration,
		&i.LocalPlayer,
		&i.DiscordID,
		&i.DiscordTag,
		&i.Username,
		&i.Avatar,
	)
	return i, err
}

const getRaidStats = `-- name: GetRaidStats :many
SELECT boss, difficulty, count(*)
FROM encounters
GROUP BY boss, difficulty
`

type GetRaidStatsRow struct {
	Boss       string
	Difficulty string
	Count      int64
}

func (q *Queries) GetRaidStats(ctx context.Context) ([]GetRaidStatsRow, error) {
	rows, err := q.db.Query(ctx, getRaidStats)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetRaidStatsRow
	for rows.Next() {
		var i GetRaidStatsRow
		if err := rows.Scan(&i.Boss, &i.Difficulty, &i.Count); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getRoles = `-- name: GetRoles :one
SELECT roles
FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetRoles(ctx context.Context, id pgtype.UUID) ([]string, error) {
	row := q.db.QueryRow(ctx, getRoles, id)
	var roles []string
	err := row.Scan(&roles)
	return roles, err
}

const getUniqueGroup = `-- name: GetUniqueGroup :one
SELECT unique_group
FROM encounters
WHERE unique_hash = $1
  AND unique_group = id
  AND (date + interval '5 minutes') >= $2
  AND (date - interval '5 minutes') <= $2
  AND (duration + 1000) >= $3
  AND (duration - 1000) <= $3
`

type GetUniqueGroupParams struct {
	UniqueHash string
	Date       pgtype.Timestamp
	Duration   int32
}

func (q *Queries) GetUniqueGroup(ctx context.Context, arg GetUniqueGroupParams) (int32, error) {
	row := q.db.QueryRow(ctx, getUniqueGroup, arg.UniqueHash, arg.Date, arg.Duration)
	var unique_group int32
	err := row.Scan(&unique_group)
	return unique_group, err
}

const getUniqueUploaders = `-- name: GetUniqueUploaders :one
SELECT COUNT(DISTINCT jsonb_object_keys(header -> 'players'))
FROM encounters
`

func (q *Queries) GetUniqueUploaders(ctx context.Context) (int64, error) {
	row := q.db.QueryRow(ctx, getUniqueUploaders)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const getUser = `-- name: GetUser :one
SELECT id, username, created_at, updated_at, access_token, discord_id, discord_tag, avatar, friends, settings, titles, roles
FROM users
WHERE discord_tag = $1
LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, discordTag string) (User, error) {
	row := q.db.QueryRow(ctx, getUser, discordTag)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.AccessToken,
		&i.DiscordID,
		&i.DiscordTag,
		&i.Avatar,
		&i.Friends,
		&i.Settings,
		&i.Titles,
		&i.Roles,
	)
	return i, err
}

const getUserByID = `-- name: GetUserByID :one
SELECT id, username, created_at, updated_at, access_token, discord_id, discord_tag, avatar, friends, settings, titles, roles
FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUserByID(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUserByID, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.AccessToken,
		&i.DiscordID,
		&i.DiscordTag,
		&i.Avatar,
		&i.Friends,
		&i.Settings,
		&i.Titles,
		&i.Roles,
	)
	return i, err
}

const getUserByToken = `-- name: GetUserByToken :one
SELECT id, username, created_at, updated_at, access_token, discord_id, discord_tag, avatar, friends, settings, titles, roles
FROM users
WHERE access_token = $1
LIMIT 1
`

func (q *Queries) GetUserByToken(ctx context.Context, accessToken pgtype.Text) (User, error) {
	row := q.db.QueryRow(ctx, getUserByToken, accessToken)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.AccessToken,
		&i.DiscordID,
		&i.DiscordTag,
		&i.Avatar,
		&i.Friends,
		&i.Settings,
		&i.Titles,
		&i.Roles,
	)
	return i, err
}

const insertEncounter = `-- name: InsertEncounter :one
INSERT
INTO encounters (uploaded_by, settings, tags, header, data, difficulty, boss, date, duration, local_player, unique_hash,
                 unique_group)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
RETURNING id
`

type InsertEncounterParams struct {
	UploadedBy  pgtype.UUID
	Settings    structs.EncounterSettings
	Tags        []string
	Header      structs.EncounterHeader
	Data        structs.EncounterData
	Difficulty  string
	Boss        string
	Date        pgtype.Timestamp
	Duration    int32
	LocalPlayer string
	UniqueHash  string
	UniqueGroup int32
}

func (q *Queries) InsertEncounter(ctx context.Context, arg InsertEncounterParams) (int32, error) {
	row := q.db.QueryRow(ctx, insertEncounter,
		arg.UploadedBy,
		arg.Settings,
		arg.Tags,
		arg.Header,
		arg.Data,
		arg.Difficulty,
		arg.Boss,
		arg.Date,
		arg.Duration,
		arg.LocalPlayer,
		arg.UniqueHash,
		arg.UniqueGroup,
	)
	var id int32
	err := row.Scan(&id)
	return id, err
}

type InsertPlayerParams struct {
	Encounter int32
	Class     string
	Name      string
	Dead      bool
	Data      structs.IndexedPlayerData
	Place     int32
}

const insertPlayerInternal = `-- name: InsertPlayerInternal :exec
INSERT
INTO players (encounter, class, name, dead, data, place)
VALUES ($1, $2, $3, $4, $5, $6)
ON CONFLICT (encounter, name)
    DO UPDATE SET data = excluded.data
`

type InsertPlayerInternalParams struct {
	Encounter int32
	Class     string
	Name      string
	Dead      bool
	Data      structs.IndexedPlayerData
	Place     int32
}

func (q *Queries) InsertPlayerInternal(ctx context.Context, arg InsertPlayerInternalParams) error {
	_, err := q.db.Exec(ctx, insertPlayerInternal,
		arg.Encounter,
		arg.Class,
		arg.Name,
		arg.Dead,
		arg.Data,
		arg.Place,
	)
	return err
}

const listEncounters = `-- name: ListEncounters :many
SELECT id
FROM encounters
ORDER BY uploaded_at DESC
`

func (q *Queries) ListEncounters(ctx context.Context) ([]int32, error) {
	rows, err := q.db.Query(ctx, listEncounters)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []int32
	for rows.Next() {
		var id int32
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		items = append(items, id)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const listRecentEncounters = `-- name: ListRecentEncounters :many
SELECT u.discord_tag,
       u.username,
       e.id,
       e.difficulty,
       e.uploaded_by,
       e.uploaded_at,
       e.settings,
       e.tags,
       e.header,
       e.boss,
       e.date,
       e.duration,
       e.local_player
FROM encounters e
         JOIN users u ON e.uploaded_by = u.id
WHERE ($1::TIMESTAMP IS NULL
    OR ($1 > e.date OR ($1::TIMESTAMP = e.date AND $2::INT < e.id)))
  AND ($3::UUID IS NULL
    OR (($4::BOOLEAN AND $3 = ANY (u.friends)) OR
        ((NOT $4::BOOLEAN AND $3 = e.uploaded_by))))
ORDER BY e.date DESC,
         e.id ASC
LIMIT 6
`

type ListRecentEncountersParams struct {
	Date    pgtype.Timestamp
	ID      pgtype.Int4
	User    pgtype.UUID
	Friends pgtype.Bool
}

type ListRecentEncountersRow struct {
	DiscordTag  string
	Username    pgtype.Text
	ID          int32
	Difficulty  string
	UploadedBy  pgtype.UUID
	UploadedAt  pgtype.Timestamp
	Settings    structs.EncounterSettings
	Tags        []string
	Header      structs.EncounterHeader
	Boss        string
	Date        pgtype.Timestamp
	Duration    int32
	LocalPlayer string
}

func (q *Queries) ListRecentEncounters(ctx context.Context, arg ListRecentEncountersParams) ([]ListRecentEncountersRow, error) {
	rows, err := q.db.Query(ctx, listRecentEncounters,
		arg.Date,
		arg.ID,
		arg.User,
		arg.Friends,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListRecentEncountersRow
	for rows.Next() {
		var i ListRecentEncountersRow
		if err := rows.Scan(
			&i.DiscordTag,
			&i.Username,
			&i.ID,
			&i.Difficulty,
			&i.UploadedBy,
			&i.UploadedAt,
			&i.Settings,
			&i.Tags,
			&i.Header,
			&i.Boss,
			&i.Date,
			&i.Duration,
			&i.LocalPlayer,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const processEncounter = `-- name: ProcessEncounter :one
UPDATE encounters
SET header      = $2,
    data        = $3,
    unique_hash = $4
WHERE id = $1
RETURNING uploaded_by
`

type ProcessEncounterParams struct {
	ID         int32
	Header     structs.EncounterHeader
	Data       structs.EncounterData
	UniqueHash string
}

func (q *Queries) ProcessEncounter(ctx context.Context, arg ProcessEncounterParams) (pgtype.UUID, error) {
	row := q.db.QueryRow(ctx, processEncounter,
		arg.ID,
		arg.Header,
		arg.Data,
		arg.UniqueHash,
	)
	var uploaded_by pgtype.UUID
	err := row.Scan(&uploaded_by)
	return uploaded_by, err
}

const setAccessToken = `-- name: SetAccessToken :exec
UPDATE users
SET access_token = $2
WHERE id = $1
`

type SetAccessTokenParams struct {
	ID          pgtype.UUID
	AccessToken pgtype.Text
}

func (q *Queries) SetAccessToken(ctx context.Context, arg SetAccessTokenParams) error {
	_, err := q.db.Exec(ctx, setAccessToken, arg.ID, arg.AccessToken)
	return err
}

const setUserRoles = `-- name: SetUserRoles :exec
UPDATE users
SET roles = $2
WHERE discord_tag = $1
`

type SetUserRolesParams struct {
	DiscordTag string
	Roles      []string
}

func (q *Queries) SetUserRoles(ctx context.Context, arg SetUserRolesParams) error {
	_, err := q.db.Exec(ctx, setUserRoles, arg.DiscordTag, arg.Roles)
	return err
}

const setUsername = `-- name: SetUsername :exec
UPDATE users
SET username = $2
WHERE id = $1
`

type SetUsernameParams struct {
	ID       pgtype.UUID
	Username pgtype.Text
}

func (q *Queries) SetUsername(ctx context.Context, arg SetUsernameParams) error {
	_, err := q.db.Exec(ctx, setUsername, arg.ID, arg.Username)
	return err
}

const updateUniqueGroup = `-- name: UpdateUniqueGroup :exec
UPDATE encounters
SET unique_group = $2
WHERE id = $1
`

type UpdateUniqueGroupParams struct {
	ID          int32
	UniqueGroup int32
}

func (q *Queries) UpdateUniqueGroup(ctx context.Context, arg UpdateUniqueGroupParams) error {
	_, err := q.db.Exec(ctx, updateUniqueGroup, arg.ID, arg.UniqueGroup)
	return err
}

const upsertEncounterGroup = `-- name: UpsertEncounterGroup :exec
INSERT
INTO grouped_encounters (group_id, uploaders)
VALUES ($1, ARRAY [$2::UUID])
ON CONFLICT (group_id)
    DO UPDATE SET uploaders = array_append(grouped_encounters.uploaders, $2)
`

type UpsertEncounterGroupParams struct {
	GroupID int32
	Column2 pgtype.UUID
}

func (q *Queries) UpsertEncounterGroup(ctx context.Context, arg UpsertEncounterGroupParams) error {
	_, err := q.db.Exec(ctx, upsertEncounterGroup, arg.GroupID, arg.Column2)
	return err
}

const upsertUser = `-- name: UpsertUser :one
INSERT
INTO users (discord_id, discord_tag, avatar, settings)
VALUES ($1, $2, $3, $4)
ON CONFLICT (discord_id)
    DO UPDATE SET discord_tag = excluded.discord_tag,
                  avatar      = excluded.avatar
RETURNING id, username, created_at, updated_at, access_token, discord_id, discord_tag, avatar, friends, settings, titles, roles
`

type UpsertUserParams struct {
	DiscordID  string
	DiscordTag string
	Avatar     pgtype.Text
	Settings   structs.UserSettings
}

func (q *Queries) UpsertUser(ctx context.Context, arg UpsertUserParams) (User, error) {
	row := q.db.QueryRow(ctx, upsertUser,
		arg.DiscordID,
		arg.DiscordTag,
		arg.Avatar,
		arg.Settings,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.Username,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.AccessToken,
		&i.DiscordID,
		&i.DiscordTag,
		&i.Avatar,
		&i.Friends,
		&i.Settings,
		&i.Titles,
		&i.Roles,
	)
	return i, err
}
