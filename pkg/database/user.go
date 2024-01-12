package database

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/0fau/logs/pkg/database/sql/structs"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
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

	var user sql.User
	if err := crdbpgx.ExecuteTx(ctx, db.Pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		qtx := db.Queries.WithTx(tx)

		row, err := qtx.GetRolesByDiscordID(ctx, discordID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return errors.Wrap(err, "getting roles")
		}

		role, err := qtx.FetchWhitelist(ctx, discordID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return errors.Wrap(err, "fetching whitelist")
			}
		} else {
			row.Roles = append(row.Roles, role)
		}

		user, err = qtx.UpsertUser(ctx, sql.UpsertUserParams{
			DiscordID:  discordID,
			DiscordTag: discordTag,
			Roles:      row.Roles,
			Avatar:     a,
			Settings: structs.UserSettings{
				SkipLanding:       false,
				LogVisibility:     "unlisted",
				ProfileVisibility: "unlisted",
			},
		})
		if err != nil {
			return errors.Wrap(err, "upserting user")
		}
		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "executing transaction")
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
