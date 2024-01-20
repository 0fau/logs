package database

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/bwmarrin/discordgo"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

func (db *DB) SaveUser(
	ctx context.Context,
	dg *discordgo.Session,
	duser *discordgo.User,
) (*sql.User, error) {
	username := duser.Username
	if duser.Discriminator != "0" {
		username += "#" + duser.Discriminator
	}

	var user sql.User
	if err := crdbpgx.ExecuteTx(ctx, db.Pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		qtx := db.Queries.WithTx(tx)

		row, err := qtx.GetRolesByDiscordID(ctx, duser.ID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return errors.Wrap(err, "getting roles")
		}

		role, err := qtx.FetchWhitelist(ctx, duser.ID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return errors.Wrap(err, "fetching whitelist")
			}
		} else {
			row.Roles = append(row.Roles, role)
		}

		user, err = qtx.UpsertUser(ctx, sql.UpsertUserParams{
			DiscordID:  duser.ID,
			DiscordTag: username,
			Roles:      row.Roles,
			Avatar:     duser.Avatar,
			Settings: structs.UserSettings{
				SkipLanding:       false,
				LogVisibility:     "unlisted",
				ProfileVisibility: "unlisted",
			},
		})
		if err != nil {
			return errors.Wrap(err, "upserting user")
		}

		if user.Avatar != duser.Avatar {
			if duser.Avatar == "" {

			}
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
