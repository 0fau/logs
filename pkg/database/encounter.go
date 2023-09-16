package database

import (
	"context"
	"fmt"
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
		LocalPlayer:      raw.LocalPlayer,
		Date:             pgtype.Timestamp{Time: date, Valid: true},
	})
	if err != nil {
		return nil, err
	}

	// buffs := []sql.InsertBuffParams{}
	skills := []sql.InsertSkillParams{}
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

		if entity.EntityType != "PLAYER" {
			continue
		}

		//if err := logBuffDamage(
		//	&buffs, enc.ID,
		//	entity.Name, entity.DamageStats.DamageDealt,
		//	entity.DamageStats.BuffedBy,
		//); err != nil {
		//	return nil, errors.Wrap(err, "logging entity buff damage")
		//}
		//
		//if err := logBuffDamage(
		//	&buffs, enc.ID,
		//	entity.Name, entity.DamageStats.DamageDealt,
		//	entity.DamageStats.DebuffedBy,
		//); err != nil {
		//	return nil, errors.Wrap(err, "logging entity debuff damage")
		//}

		for idstr, skill := range entity.Skills {
			id, err := strconv.Atoi(idstr)
			if err != nil {
				return nil, errors.Wrap(err, "parsing skill id")
			}

			skills = append(skills, sql.InsertSkillParams{
				Encounter:   enc.ID,
				Player:      entity.Name,
				SkillID:     int32(id),
				Casts:       skill.Casts,
				Crits:       skill.Crits,
				Dps:         skill.DPS,
				Hits:        skill.Hits,
				MaxDamage:   skill.MaxDamage,
				TotalDamage: skill.TotalDamage,
				Name:        skill.Name,
			})
		}
	}
	if _, err := qtx.InsertSkill(ctx, skills); err != nil {
		return nil, errors.Wrap(err, "inserting skills")
	}
	//if _, err := qtx.InsertBuff(ctx, buffs); err != nil {
	//	return nil, errors.Wrap(err, "inserting buffs")
	//}

	if err := tx.Commit(ctx); err != nil {
		return nil, err
	}

	return &enc, nil
}

func logBuffDamage(
	arr *[]sql.InsertBuffParams,
	enc int32, player string, totalDamage int64,
	buffs map[string]int64,
) error {
	for id, damage := range buffs {
		buff, err := strconv.Atoi(id)
		if err != nil {
			return errors.Wrap(err, "parsing skill id")
		}

		percent := float64(damage) / float64(totalDamage) * 100
		numeric := new(pgtype.Numeric)
		if err := numeric.Scan(fmt.Sprintf("%.2f", percent)); err != nil {
			return errors.Wrap(err, "scanning damage percent into numeric")
		}

		*arr = append(*arr, sql.InsertBuffParams{
			Encounter: enc,
			Player:    player,
			BuffID:    int32(buff),
			Damage:    damage,
			Percent:   *numeric,
		})
	}
	return nil
}

func (db *DB) RecentEncounters(ctx context.Context) ([]*sql.Encounter, error) {
	encounters, err := db.queries.ListRecentEncounters(ctx)
	if err != nil {
		return nil, err
	}
	ret := make([]*sql.Encounter, len(encounters))
	for i := 0; i < len(encounters); i++ {
		ret[i] = &encounters[i]
	}
	return ret, nil
}

func (db *DB) ListEntities(ctx context.Context, enc int32) ([]*sql.Entity, error) {
	entities, err := db.queries.GetEntities(ctx, enc)
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
	skills, err := db.queries.GetSkills(ctx, enc)
	if err != nil {
		return nil, errors.Wrapf(err, "getting skills for encounter %d", enc)
	}
	ret := make([]*sql.Skill, len(skills))
	for i := 0; i < len(skills); i++ {
		ret[i] = &skills[i]
	}
	return ret, nil
}
