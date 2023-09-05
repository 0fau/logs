package database

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/0fau/logs/pkg/process"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

func (db *DB) SaveEncounter(uuid pgtype.UUID, raw *process.RawEncounter) (*sql.Encounter, error) {
	ctx := context.Background()
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := db.queries.WithTx(tx)

	date := time.UnixMilli(raw.FightStart)
	enc, err := qtx.InsertEncounter(ctx, sql.InsertEncounterParams{
		UploadedBy:       uuid,
		Visibility:       "unlisted",
		Raid:             raw.CurrentBossName,
		TotalDamageDealt: raw.DamageStats.TotalDamageDealt,
		Cleared:          raw.DamageStats.Misc.Cleared,
		Duration:         raw.Duration,
		Date:             pgtype.Timestamp{Time: date, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	for _, entity := range raw.Entities {
		_, err := qtx.InsertEntity(ctx, sql.InsertEntityParams{
			Encounter: enc.ID,
			Class:     entity.Class,
			Enttype:   entity.EntityType,
			Name:      entity.Name,
			Damage:    entity.DamageStats.DamageDealt,
			Dps:       entity.DamageStats.DPS,
		})
		if err != nil {
			return nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &enc, nil
}

func (db *DB) RecentEncounters() {

}
