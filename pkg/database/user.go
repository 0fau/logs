package database

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/jackc/pgx/v5/pgtype"
)

func (db *DB) SaveUser(
	ctx context.Context,
	discordID, discordName string,
) (*sql.User, error) {
	user, err := db.queries.UpsertUser(ctx, sql.UpsertUserParams{
		DiscordID:   discordID,
		DiscordName: discordName,
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

	return db.queries.SetAccessToken(ctx, sql.SetAccessTokenParams{
		ID:          uuid,
		AccessToken: pgtext(token),
	})
}

func (db *DB) UserByAccessToken(
	ctx context.Context,
	token string,
) (*sql.User, error) {
	user, err := db.queries.GetUserByToken(ctx, pgtext(token))
	if err != nil {
		return nil, err
	}
	return &user, err
}
