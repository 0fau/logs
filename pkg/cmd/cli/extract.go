package main

import (
	"database/sql"
	"fmt"
	"github.com/0fau/logs/pkg/process/meter"
	"github.com/goccy/go-json"
	"github.com/spf13/cobra"
	"log"
)

func extract() *cobra.Command {
	return &cobra.Command{
		Use:  "extract",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			db, err := sql.Open("sqlite3", "./sample/encounters.db")
			if err != nil {
				log.Fatal(err)
			}
			defer db.Close()

			encID := args[0]

			statement := `
    SELECT
       fight_start,
       local_player,
       current_boss,
       duration,
       last_combat_packet,
       total_damage_dealt,
       buffs,
       debuffs,
       misc,
       difficulty
    FROM encounter
    WHERE id = ?
    LIMIT 1;
`
			rows, err := db.Query(statement, encID)
			if err != nil {
				log.Fatal(err)
			}
			defer rows.Close()

			if !rows.Next() {
				return sql.ErrNoRows
			}
			var enc meter.Encounter

			var (
				buffs      string
				debuffs    string
				misc       string
				difficulty sql.NullString
			)

			if err := rows.Scan(
				&enc.FightStart,
				&enc.LocalPlayer,
				&enc.CurrentBossName,
				&enc.Duration,
				&enc.End,
				&enc.DamageStats.TotalDamageDealt,
				&buffs,
				&debuffs,
				&misc,
				&difficulty,
			); err != nil {
				return err
			}

			enc.Difficulty = difficulty.String

			for _, unmarshal := range []struct {
				in  string
				out interface{}
			}{
				{misc, &enc.DamageStats.Misc},
				{buffs, &enc.DamageStats.Buffs},
				{debuffs, &enc.DamageStats.Debuffs},
			} {
				if err := json.Unmarshal([]byte(unmarshal.in), unmarshal.out); err != nil {
					return err
				}
			}

			entStatement := `
    SELECT name,
        class_id,
        class,
        gear_score,
        is_dead,
        skills,
        damage_stats,
        skill_stats,
        entity_type
    FROM entity
    WHERE encounter_id = ?;
`

			entRows, err := db.Query(entStatement, encID)
			if err != nil {
				return err
			}
			defer entRows.Close()

			enc.Entities = make(map[string]meter.Entity)
			for entRows.Next() {
				var ent meter.Entity

				var (
					skills      string
					damageStats string
					skillStats  string
				)

				if err := entRows.Scan(
					&ent.Name,
					&ent.ClassId,
					&ent.Class,
					&ent.GearScore,
					&ent.Dead,
					&skills,
					&damageStats,
					&skillStats,
					&ent.EntityType,
				); err != nil {
					return err
				}

				for _, unmarshal := range []struct {
					in  string
					out interface{}
				}{
					{skills, &ent.Skills},
					{damageStats, &ent.DamageStats},
					{skillStats, &ent.SkillStats},
				} {
					if err := json.Unmarshal([]byte(unmarshal.in), unmarshal.out); err != nil {
						return err
					}
				}

				enc.Entities[ent.Name] = ent
			}

			raw, _ := json.MarshalIndent(enc, "", "  ")
			fmt.Println(string(raw))

			err = rows.Err()
			if err != nil {
				return err
			}

			return nil
		},
	}
}
