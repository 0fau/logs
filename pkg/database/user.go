package database

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5/pgtype"
)

func (db *DB) SaveUser(
	ctx context.Context,
	discordID, discordTag string, avatar *string,
) (*sql.User, error) {
	var a pgtype.Text
	if avatar != nil {
		if err := a.Scan(*avatar); err != nil {
			return nil, errors.Wrap(err, "scanning avatar pgtype.Text")
		}
	}

	user, err := db.Queries.UpsertUser(ctx, sql.UpsertUserParams{
		DiscordID:  discordID,
		DiscordTag: discordTag,
		Avatar:     a,
		Settings: structs.UserSettings{
			SkipLanding:       false,
			LogVisibility:     "unlisted",
			ProfileVisibility: "unlisted",
		},
	})
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (db *DB) SetUserAccessToken(
	ctx context.Context,
	id string, token string,
) error {
	var uuid pgtype.UUID
	if err := uuid.Scan(id); err != nil {
		return err
	}

	at := pgtype.Text{}
	if token != "" {
		if err := at.Scan(token); err != nil {
			return errors.Wrap(err, "scanning pgtype.Text access token")
		}
	}

	return db.Queries.SetAccessToken(ctx, sql.SetAccessTokenParams{
		ID:          uuid,
		AccessToken: at,
	})
}

func (db *DB) UserByAccessToken(
	ctx context.Context,
	token string,
) (*sql.User, error) {
	user, err := db.Queries.GetUserByToken(ctx, pgtext(token))
	if err != nil {
		return nil, err
	}
	return &user, err
}
