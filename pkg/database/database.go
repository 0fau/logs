package database

import (
	"context"
	"errors"
	"github.com/jackc/pgx/v5/pgxpool"
	"net/url"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/cockroachdb"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func Migrate(dbURL string) error {
	u, err := url.Parse(dbURL)
	if err != nil {
		return err
	}
	u.Scheme = "cockroachdb"

	m, err := migrate.New("file://migrations", u.String())
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}

	return nil
}

func NewPool(ctx context.Context, dbURL string) (*pgxpool.Pool, error) {
	config, err := pgxpool.ParseConfig(dbURL)
	if err != nil {
		return nil, err
	}

	return pgxpool.NewWithConfig(ctx, config)
}
