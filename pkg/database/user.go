package database

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5/pgtype"
)

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
