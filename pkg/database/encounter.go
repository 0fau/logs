package database

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

func (db *DB) RecentEncounters(ctx context.Context, user string, id int32, date time.Time) ([]sql.ListRecentEncountersRow, error) {
	args := sql.ListRecentEncountersParams{
		Date: pgtype.Timestamp{
			Time:  date,
			Valid: !date.IsZero(),
		},
		ID: pgtype.Int4{
			Int32: id,
			Valid: id != 0,
		},
		User: pgtype.UUID{},
	}

	if user != "" {
		err := args.User.Scan(user)
		if err != nil {
			return nil, errors.Wrap(err, "scanning user uuid")
		}
	}

	return db.Queries.ListRecentEncounters(ctx, args)
}
