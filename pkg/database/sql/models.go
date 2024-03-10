// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.23.0

package sql

import (
	structs "github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/jackc/pgx/v5/pgtype"
)

type Encounter struct {
	ID          int32
	UploadedBy  pgtype.UUID
	UploadedAt  pgtype.Timestamp
	Settings    structs.EncounterSettings
	Thumbnail   bool
	Tags        []string
	Header      structs.EncounterHeader
	Data        structs.EncounterData
	UniqueHash  string
	UniqueGroup int32
	Visibility  *structs.EncounterVisibility
	Boss        string
	Difficulty  string
	Date        pgtype.Timestamp
	Duration    int32
	Version     int32
	LocalPlayer string
}

type Friend struct {
	User1 pgtype.UUID
	User2 pgtype.UUID
	Date  pgtype.Timestamp
}

type FriendRequest struct {
	User1 string
	User2 string
	Date  pgtype.Timestamp
}

type GroupedEncounter struct {
	GroupID   int32
	Uploaders []pgtype.UUID
}

type Player struct {
	Encounter  int32
	Boss       string
	Difficulty string
	Class      string
	Name       string
	Dead       bool
	Dps        int64
	GearScore  float64
	Place      int32
}

type Roster struct {
	UserID    pgtype.UUID
	Character string
	Class     string
	GearScore float64
}

type User struct {
	ID            pgtype.UUID
	Username      *string
	CreatedAt     pgtype.Timestamp
	UpdatedAt     pgtype.Timestamp
	AccessToken   *string
	DiscordID     string
	DiscordTag    string
	Avatar        string
	Friends       []pgtype.UUID
	Settings      structs.UserSettings
	LogVisibility *structs.EncounterVisibility
	Titles        []string
	Roles         []string
}

type Whitelist struct {
	Discord string
	Role    string
}
