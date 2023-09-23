package database

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/0fau/logs/pkg/process/meter"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
	"time"
)

func (db *DB) SaveEncounter(
	ctx context.Context,
	uuid pgtype.UUID,
	raw *meter.Encounter,
) (*sql.Encounter, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback(ctx)
	qtx := db.Queries.WithTx(tx)

	date := time.UnixMilli(raw.FightStart)
	enc, err := qtx.InsertEncounter(ctx, sql.InsertEncounterParams{
		UploadedBy: uuid,
		Visibility: "unlisted",
		Raid:       raw.CurrentBossName,
		Damage:     raw.DamageStats.TotalDamageDealt,
		Cleared:    raw.DamageStats.Misc.Cleared,
		Duration:   raw.Duration,
		Fields: meter.StoredEncounterFields{
			Buffs:     raw.DamageStats.Buffs,
			Debuffs:   raw.DamageStats.Debuffs,
			PartyInfo: raw.DamageStats.Misc.PartyInfo,
			HPLog:     raw.DamageStats.Misc.HPLog,
		},
		LocalPlayer: raw.LocalPlayer,
		Date:        pgtype.Timestamp{Time: date, Valid: true},
	})
	if err != nil {
		return nil, errors.Wrap(err, "inserting encounter")
	}

	var skills []sql.InsertSkillParams
	for _, entity := range raw.Entities {
		_, err := qtx.InsertEntity(ctx, sql.InsertEntityParams{
			Encounter: enc.ID,
			Class:     entity.Class,
			Enttype:   entity.EntityType,
			Name:      entity.Name,
			Damage:    entity.DamageStats.DamageDealt,
			Dps:       entity.DamageStats.DPS,
			Dead:      entity.Dead,
			Fields: meter.StoredEntityFields{
				Buffed:         entity.DamageStats.BuffedBy,
				Debuffed:       entity.DamageStats.DebuffedBy,
				BuffedDamage:   entity.DamageStats.BuffedDamage,
				DebuffedDamage: entity.DamageStats.DebuffedDamage,
				FADamage:       entity.DamageStats.FADamage,
				BADamage:       entity.DamageStats.BADamage,
				DeathTime:      entity.DamageStats.DeathTime,
				DPSAverage:     entity.DamageStats.DPSAverage,
				DPSRolling:     entity.DamageStats.DPSRolling,
			},
		})
		if err != nil {
			return nil, errors.Wrap(err, "inserting entity")
		}

		if entity.EntityType != "PLAYER" {
			continue
		}

		for idstr, skill := range entity.Skills {
			id, err := strconv.Atoi(idstr)
			if err != nil {
				return nil, errors.Wrap(err, "parsing skill id")
			}

			skills = append(skills, sql.InsertSkillParams{
				Encounter: enc.ID,
				Player:    entity.Name,
				SkillID:   int32(id),
				Tripods:   skill.TripodIndex,
				Fields: meter.StoredSkillFields{
					Casts:        skill.Casts,
					CastLog:      skill.CastLog,
					Crits:        skill.Crits,
					Hits:         skill.Hits,
					Buffed:       skill.BuffedBy,
					Debuffed:     skill.DebuffedBy,
					MaxDamage:    skill.MaxDamage,
					FADamage:     skill.FADamage,
					BADamage:     skill.BADamage,
					TripodLevels: skill.TripodLevel,
				},
				Dps:    skill.DPS,
				Damage: skill.TotalDamage,
				Name:   skill.Name,
			})
		}
	}
	if _, err := qtx.InsertSkill(ctx, skills); err != nil {
		return nil, errors.Wrap(err, "inserting skills")
	}

	if err := tx.Commit(ctx); err != nil {
		return nil, errors.Wrap(err, "committing transaction")
	}

	return &enc, nil
}

func (db *DB) RecentEncounters(ctx context.Context) ([]*sql.ListRecentEncountersRow, error) {
	encounters, err := db.Queries.ListRecentEncounters(ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*sql.ListRecentEncountersRow, len(encounters))
	for i := 0; i < len(encounters); i++ {
		ret[i] = &encounters[i]
	}
	return ret, nil
}

func (db *DB) ListEntities(ctx context.Context, enc int32) ([]*sql.Entity, error) {
	entities, err := db.Queries.GetEntities(ctx, enc)
	if err != nil {
		return nil, err
	}
	ret := make([]*sql.Entity, len(entities))
	for i := 0; i < len(entities); i++ {
		ret[i] = &entities[i]
	}
	return ret, nil
}

func (db *DB) ListSkills(ctx context.Context, enc int32) ([]*sql.Skill, error) {
	skills, err := db.Queries.GetSkills(ctx, enc)
	if err != nil {
		return nil, errors.Wrapf(err, "getting skills for encounter %d", enc)
	}
	ret := make([]*sql.Skill, len(skills))
	for i := 0; i < len(skills); i++ {
		ret[i] = &skills[i]
	}
	return ret, nil
}
