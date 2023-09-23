// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0
// source: query.sql

package sql

import (
	"context"

	meter "github.com/0fau/logs/pkg/process/meter"
	"github.com/jackc/pgx/v5/pgtype"
)

const deleteUser = `-- name: DeleteUser :exec
DELETE
FROM users
WHERE id = $1
`

func (q *Queries) DeleteUser(ctx context.Context, id pgtype.UUID) error {
	_, err := q.db.Exec(ctx, deleteUser, id)
	return err
}

const getEncounter = `-- name: GetEncounter :one
SELECT id, uploaded_by, visibility, title, description, raid, date, duration, damage, fields, cleared, uploaded_at, tags, local_player
FROM encounters
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetEncounter(ctx context.Context, id int32) (Encounter, error) {
	row := q.db.QueryRow(ctx, getEncounter, id)
	var i Encounter
	err := row.Scan(
		&i.ID,
		&i.UploadedBy,
		&i.Visibility,
		&i.Title,
		&i.Description,
		&i.Raid,
		&i.Date,
		&i.Duration,
		&i.Damage,
		&i.Fields,
		&i.Cleared,
		&i.UploadedAt,
		&i.Tags,
		&i.LocalPlayer,
	)
	return i, err
}

const getEntities = `-- name: GetEntities :many
SELECT encounter, enttype, name, class, damage, dps, dead, fields
FROM entities
WHERE encounter = $1
`

func (q *Queries) GetEntities(ctx context.Context, encounter int32) ([]Entity, error) {
	rows, err := q.db.Query(ctx, getEntities, encounter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entity
	for rows.Next() {
		var i Entity
		if err := rows.Scan(
			&i.Encounter,
			&i.Enttype,
			&i.Name,
			&i.Class,
			&i.Damage,
			&i.Dps,
			&i.Dead,
			&i.Fields,
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

const getFields = `-- name: GetFields :one
SELECT fields
FROM encounters
WHERE id = $1
`

func (q *Queries) GetFields(ctx context.Context, id int32) (meter.StoredEncounterFields, error) {
	row := q.db.QueryRow(ctx, getFields, id)
	var fields meter.StoredEncounterFields
	err := row.Scan(&fields)
	return fields, err
}

const getSkills = `-- name: GetSkills :many
SELECT encounter, player, skill_id, name, dps, damage, tripods, fields
FROM skills
WHERE encounter = $1
`

func (q *Queries) GetSkills(ctx context.Context, encounter int32) ([]Skill, error) {
	rows, err := q.db.Query(ctx, getSkills, encounter)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Skill
	for rows.Next() {
		var i Skill
		if err := rows.Scan(
			&i.Encounter,
			&i.Player,
			&i.SkillID,
			&i.Name,
			&i.Dps,
			&i.Damage,
			&i.Tripods,
			&i.Fields,
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

const getUser = `-- name: GetUser :one
SELECT id, discord_id, discord_name, access_token, roles, created_at, updated_at
FROM users
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetUser(ctx context.Context, id pgtype.UUID) (User, error) {
	row := q.db.QueryRow(ctx, getUser, id)
	var i User
	err := row.Scan(
		&i.ID,
		&i.DiscordID,
		&i.DiscordName,
		&i.AccessToken,
		&i.Roles,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const getUserByToken = `-- name: GetUserByToken :one
SELECT id, discord_id, discord_name, access_token, roles, created_at, updated_at
FROM users
WHERE access_token = $1
LIMIT 1
`

func (q *Queries) GetUserByToken(ctx context.Context, accessToken pgtype.Text) (User, error) {
	row := q.db.QueryRow(ctx, getUserByToken, accessToken)
	var i User
	err := row.Scan(
		&i.ID,
		&i.DiscordID,
		&i.DiscordName,
		&i.AccessToken,
		&i.Roles,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}

const insertEncounter = `-- name: InsertEncounter :one
INSERT
INTO encounters (uploaded_by, raid, date, visibility, duration, damage, cleared, local_player, fields)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
RETURNING id, uploaded_by, visibility, title, description, raid, date, duration, damage, fields, cleared, uploaded_at, tags, local_player
`

type InsertEncounterParams struct {
	UploadedBy  pgtype.UUID
	Raid        string
	Date        pgtype.Timestamp
	Visibility  string
	Duration    int32
	Damage      int64
	Cleared     bool
	LocalPlayer string
	Fields      meter.StoredEncounterFields
}

func (q *Queries) InsertEncounter(ctx context.Context, arg InsertEncounterParams) (Encounter, error) {
	row := q.db.QueryRow(ctx, insertEncounter,
		arg.UploadedBy,
		arg.Raid,
		arg.Date,
		arg.Visibility,
		arg.Duration,
		arg.Damage,
		arg.Cleared,
		arg.LocalPlayer,
		arg.Fields,
	)
	var i Encounter
	err := row.Scan(
		&i.ID,
		&i.UploadedBy,
		&i.Visibility,
		&i.Title,
		&i.Description,
		&i.Raid,
		&i.Date,
		&i.Duration,
		&i.Damage,
		&i.Fields,
		&i.Cleared,
		&i.UploadedAt,
		&i.Tags,
		&i.LocalPlayer,
	)
	return i, err
}

const insertEntity = `-- name: InsertEntity :one
INSERT
INTO entities (encounter, class, enttype, name, damage, dps, dead, fields)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING encounter, enttype, name, class, damage, dps, dead, fields
`

type InsertEntityParams struct {
	Encounter int32
	Class     string
	Enttype   string
	Name      string
	Damage    int64
	Dps       int64
	Dead      bool
	Fields    meter.StoredEntityFields
}

func (q *Queries) InsertEntity(ctx context.Context, arg InsertEntityParams) (Entity, error) {
	row := q.db.QueryRow(ctx, insertEntity,
		arg.Encounter,
		arg.Class,
		arg.Enttype,
		arg.Name,
		arg.Damage,
		arg.Dps,
		arg.Dead,
		arg.Fields,
	)
	var i Entity
	err := row.Scan(
		&i.Encounter,
		&i.Enttype,
		&i.Name,
		&i.Class,
		&i.Damage,
		&i.Dps,
		&i.Dead,
		&i.Fields,
	)
	return i, err
}

type InsertSkillParams struct {
	Encounter int32
	Player    string
	SkillID   int32
	Dps       int64
	Damage    int64
	Name      string
	Tripods   meter.TripodRows
	Fields    meter.StoredSkillFields
}

const listRecentEncounters = `-- name: ListRecentEncounters :many
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
LIMIT 5
`

type ListRecentEncountersRow struct {
	ID          int32
	UploadedBy  pgtype.UUID
	Raid        string
	Date        pgtype.Timestamp
	Visibility  string
	Duration    int32
	Damage      int64
	Cleared     bool
	LocalPlayer string
}

func (q *Queries) ListRecentEncounters(ctx context.Context) ([]ListRecentEncountersRow, error) {
	rows, err := q.db.Query(ctx, listRecentEncounters)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ListRecentEncountersRow
	for rows.Next() {
		var i ListRecentEncountersRow
		if err := rows.Scan(
			&i.ID,
			&i.UploadedBy,
			&i.Raid,
			&i.Date,
			&i.Visibility,
			&i.Duration,
			&i.Damage,
			&i.Cleared,
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

const upsertUser = `-- name: UpsertUser :one
INSERT
INTO users (discord_id, discord_name)
VALUES ($1, $2)
ON CONFLICT (discord_id)
    DO UPDATE SET discord_name = excluded.discord_name
RETURNING id, discord_id, discord_name, access_token, roles, created_at, updated_at
`

type UpsertUserParams struct {
	DiscordID   string
	DiscordName string
}

func (q *Queries) UpsertUser(ctx context.Context, arg UpsertUserParams) (User, error) {
	row := q.db.QueryRow(ctx, upsertUser, arg.DiscordID, arg.DiscordName)
	var i User
	err := row.Scan(
		&i.ID,
		&i.DiscordID,
		&i.DiscordName,
		&i.AccessToken,
		&i.Roles,
		&i.CreatedAt,
		&i.UpdatedAt,
	)
	return i, err
}
