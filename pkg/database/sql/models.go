// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.20.0

package sql

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type User struct {
	ID          pgtype.UUID
	DiscordID   pgtype.Text
	DiscordName pgtype.Text
	AccessToken pgtype.Text
	Roles       []string
	CreatedAt   pgtype.Timestamp
	UpdatedAt   pgtype.Timestamp
}
