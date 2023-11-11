package database

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5/pgtype"
	"time"
)

//func (db *DB) SaveEncounter(
//	ctx context.Context,
//	uuid pgtype.UUID,
//	raw *meter.Encounter,
//) (int32, error) {
//	tx, err := db.Pool.Begin(ctx)
//	if err != nil {
//		return 0, err
//	}
//	defer tx.Rollback(ctx)
//	qtx := db.Queries.WithTx(tx)
//
//	date := time.UnixMilli(raw.FightStart)
//	enc, err := qtx.InsertEncounter(ctx, sql.InsertEncounterParams{
//		UploadedBy: uuid,
//		Visibility: "unlisted",
//		Boss:       raw.CurrentBossName,
//		Damage:     raw.DamageStats.TotalDamageDealt,
//		Cleared:    raw.DamageStats.Misc.Cleared,
//		Duration:   raw.Duration,
//		Fields: meter.StoredEncounterFields{
//			Buffs:     raw.DamageStats.Buffs,
//			Debuffs:   raw.DamageStats.Debuffs,
//			PartyInfo: raw.DamageStats.Misc.PartyInfo,
//			HPLog:     raw.DamageStats.Misc.HPLog,
//		},
//		LocalPlayer: raw.LocalPlayer,
//		Date:        pgtype.Timestamp{Time: date, Valid: true},
//	})
//	if err != nil {
//		return 0, errors.Wrap(err, "inserting encounter")
//	}
//
//	var skills []sql.InsertSkillParams
//	for _, entity := range raw.Entities {
//		_, err := qtx.InsertEntity(ctx, sql.InsertEntityParams{
//			Encounter: enc,
//			Class:     entity.Class,
//			Enttype:   entity.EntityType,
//			Name:      entity.Name,
//			Damage:    entity.DamageStats.Damage,
//			Dps:       entity.DamageStats.DPS,
//			Dead:      entity.Dead,
//			Fields: meter.StoredEntityFields{
//				BuffedBy:   entity.DamageStats.BuffedBy,
//				DebuffedBy: entity.DamageStats.DebuffedBy,
//				BuffedBy:     entity.DamageStats.BuffedBy,
//				DebuffedBy:   entity.DamageStats.DebuffedBy,
//				FA:         entity.SkillStats.FA,
//				BA:         entity.SkillStats.BA,
//				FADamage:   entity.DamageStats.FADamage,
//				BADamage:   entity.DamageStats.BADamage,
//				DeathTime:  entity.DamageStats.DeathTime,
//				DPSAverage: entity.DamageStats.DPSAverage,
//				DPSRolling: entity.DamageStats.DPSRolling,
//			},
//		})
//		if err != nil {
//			return 0, errors.Wrap(err, "inserting entity")
//		}
//
//		if entity.EntityType != "PLAYER" {
//			continue
//		}
//
//		for idstr, skill := range entity.Skills {
//			id, err := strconv.Atoi(idstr)
//			if err != nil {
//				return 0, errors.Wrap(err, "parsing skill id")
//			}
//
//			skills = append(skills, sql.InsertSkillParams{
//				Encounter: enc,
//				Player:    entity.Name,
//				SkillID:   int32(id),
//				Tripods:   skill.TripodIndex,
//				Fields: meter.StoredSkillFields{
//					Casts:        skill.Casts,
//					CastLog:      skill.CastLog,
//					Crits:        skill.Crits,
//					Hits:         skill.Hits,
//					Icon:         skill.Icon,
//					BuffedBy:       skill.BuffedBy,
//					DebuffedBy:     skill.DebuffedBy,
//					MaxDamage:    skill.Max,
//					FADamage:     skill.FADamage,
//					BADamage:     skill.BADamage,
//					TripodLevels: skill.TripodLevel,
//				},
//				Dps:    skill.DPS,
//				Damage: skill.Damage,
//				Name:   skill.Name,
//			})
//		}
//	}
//	if _, err := qtx.InsertSkill(ctx, skills); err != nil {
//		return 0, errors.Wrap(err, "inserting skills")
//	}
//
//	if err := tx.Commit(ctx); err != nil {
//		return 0, errors.Wrap(err, "committing transaction")
//	}
//
//	return enc, nil
//}

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

func (db *DB) ListEntities(ctx context.Context, enc int32) ([]*sql.Player, error) {
	entities, err := db.Queries.GetEntities(ctx, enc)
	if err != nil {
		return nil, err
	}
	ret := make([]*sql.Player, len(entities))
	for i := 0; i < len(entities); i++ {
		ret[i] = &entities[i]
	}
	return ret, nil
}

func (db *DB) ListSkills(ctx context.Context, enc int32) ([]*sql.Skill, error) {
	skills, err := db.Queries.GetSkills(ctx, enc)
	if err != nil {
		return nil, errors.Wrapf(err, "getting skills for [encounter] %d", enc)
	}
	ret := make([]*sql.Skill, len(skills))
	for i := 0; i < len(skills); i++ {
		ret[i] = &skills[i]
	}
	return ret, nil
}
