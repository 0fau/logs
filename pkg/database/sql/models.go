// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package sql

import (
	structs "github.com/0fau/logs/pkg/database/sql/structs"
	meter "github.com/0fau/logs/pkg/process/meter"
	"github.com/jackc/pgx/v5/pgtype"
)

type Encounter struct {
	ID          int32
	UploadedBy  pgtype.UUID
	UploadedAt  pgtype.Timestamp
	Settings    structs.EncounterSettings
	Tags        []string
	Header      structs.EncounterHeader
	Data        structs.EncounterData
	Boss        string
	Date        pgtype.Timestamp
	Duration    int32
	LocalPlayer string
}

type Player struct {
	Encounter int32
	Enttype   string
	Name      string
	Class     string
	Damage    int64
	Dps       int64
	Dead      bool
	Fields    []byte
}

type Skill struct {
	Encounter int32
	Player    string
	SkillID   int32
	Name      string
	Dps       int64
	Damage    int64
	Tripods   meter.TripodRows
	Fields    meter.StoredSkillFields
}

type User struct {
	ID          pgtype.UUID
	Username    pgtype.Text
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
	AccessToken pgtype.Text
	DiscordID   string
	DiscordTag  string
	Avatar      pgtype.Text
	Friends     []pgtype.UUID
	Settings    structs.UserSettings
	Titles      []string
	Roles       []string
}
