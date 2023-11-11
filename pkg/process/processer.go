package process

import (
	"cmp"
	"context"
	"fmt"
	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/0fau/logs/pkg/process/meter"
	"github.com/0fau/logs/pkg/s3"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5/pgtype"
	"math"
	"slices"
	"strconv"
	"time"
)

type Processor struct {
	skills map[string]SkillData
	buffs  map[string]BuffData

	db *database.DB
	s3 *s3.Client
}

func NewLogProcessor(db *database.DB, s3 *s3.Client) *Processor {
	return &Processor{
		db: db,
		s3: s3,
	}
}

func (p *Processor) Initialize() error {
	if err := p.loadMeterData(); err != nil {
		return err
	}

	return nil
}

type Encounter struct {
	raw *meter.Encounter

	Header structs.EncounterHeader
	Data   structs.EncounterData
}

func (p *Processor) Lint(enc *meter.Encounter) error {
	return nil
}

func (p *Processor) Process(raw *meter.Encounter) (*Encounter, error) {
	enc := &Encounter{raw: raw}
	header, err := enc.processHeader()
	if err != nil {
		return nil, errors.Wrap(err, "processing encounter header")
	}
	enc.Header = header

	data, err := enc.processData()
	if err != nil {
		return nil, errors.Wrap(err, "processing encounter data")
	}
	enc.Data = data

	enc.highlight()
	return enc, nil
}

func (enc *Encounter) processHeader() (structs.EncounterHeader, error) {
	header := structs.EncounterHeader{
		Players: make(map[string]structs.PlayerHeader),
		Parties: make([][]string, len(enc.raw.DamageStats.Misc.PartyInfo)),
		Damage:  enc.raw.DamageStats.TotalDamageDealt,
		Cleared: enc.raw.DamageStats.Misc.Cleared,
	}

	for party, players := range enc.raw.DamageStats.Misc.PartyInfo {
		num, err := strconv.Atoi(party)
		if err != nil {
			return structs.EncounterHeader{}, errors.Wrapf(err, "converting party %s to number", party)
		}

		slices.SortFunc(players, func(a, b string) int {
			return cmp.Compare(
				enc.raw.Entities[b].DamageStats.Damage,
				enc.raw.Entities[a].DamageStats.Damage,
			)
		})
		header.Parties[num] = players
	}

	for _, entity := range enc.raw.Entities {
		if entity.EntityType == "ESTHER" {
			continue
		}

		header.Players[entity.Name] = structs.PlayerHeader{
			Name:   entity.Name,
			Class:  entity.Class,
			Damage: entity.DamageStats.Damage,
			DPS:    entity.DamageStats.DPS,
			Alive:  !entity.Dead,
		}
	}
	return header, nil
}

func round(dec float64) string {
	if math.IsNaN(dec) {
		return "0.0"
	}

	ratio := math.Pow(10, float64(1))
	dec = math.Round(dec*ratio) / ratio
	return fmt.Sprintf("%.1f", dec)
}

type Synergies map[string]Synergy

type Synergy struct {
	Name  string
	Buffs map[string]struct{}
}

func (s Synergies) Collect(group, buff string) {
	syn, ok := s[group]
	if !ok {
		syn.Name = group
		syn.Buffs = make(map[string]struct{})
	}
	syn.Buffs[buff] = struct{}{}
	s[group] = syn
}

func (s Synergies) Serialize() []structs.Synergy {
	arr := make([]structs.Synergy, 0, len(s))
	for _, syn := range s {
		ssyn := structs.Synergy{
			Name:  syn.Name,
			Buffs: make([]string, 0, len(syn.Buffs)),
		}
		for buff := range syn.Buffs {
			ssyn.Buffs = append(ssyn.Buffs, buff)
		}
		arr = append(arr, ssyn)
	}
	return arr
}

func (enc *Encounter) processData() (structs.EncounterData, error) {
	data := structs.EncounterData{
		Players:      make(map[string]structs.PlayerData),
		Synergies:    make([][]structs.Synergy, len(enc.Header.Parties)),
		BuffCatalog:  make(map[string]structs.BuffInfo),
		SkillCatalog: make(map[string]structs.SkillInfo),
	}

	parties := map[string]int{}
	for party, players := range enc.Header.Parties {
		for _, player := range players {
			parties[player] = party
		}
	}

	partyBuffs := make([]Synergies, len(enc.Header.Parties))
	for i := range partyBuffs {
		partyBuffs[i] = make(Synergies)
	}
	for name, entity := range enc.raw.Entities {
		if entity.EntityType != "PLAYER" {
			continue
		}

		data.Players[name] = enc.processPlayer(entity)
		for gname, group := range data.Players[name].Synergy {

			for buff := range group.Buffs {
				if _, ok := data.BuffCatalog[buff]; ok {
					continue
				}

				partyBuffs[parties[name]].Collect(gname, buff)

				info, ok := enc.raw.DamageStats.Buffs[buff]
				if !ok {
					info, ok = enc.raw.DamageStats.Debuffs[buff]
					if !ok {
						continue
					}
				}

				binfo := structs.BuffInfo{
					Name:        info.Source.Name,
					Icon:        info.Source.Icon,
					Description: info.Source.Description,
					Category:    info.Category,
					Set:         info.Source.SetName,
				}
				if info.Source.Skill != nil {
					binfo.Skill = &structs.BuffSkill{
						Class:       int(info.Source.Skill.ClassID),
						Description: info.Source.Skill.Description,
						Name:        info.Source.Skill.Name,
						Icon:        info.Source.Skill.Icon,
						ID:          int(info.Source.Skill.ID),
					}
				}
				data.BuffCatalog[buff] = binfo
			}
		}

		for skill := range data.Players[name].SkillDamage {
			info := enc.raw.Entities[name].Skills[skill]
			data.SkillCatalog[skill] = structs.SkillInfo{
				Name: info.Name,
				Icon: info.Icon,
			}
		}

	}

	for i, groups := range partyBuffs {
		data.Synergies[i] = groups.Serialize()
	}

	return data, nil
}

func (enc *Encounter) processPlayer(entity meter.Entity) structs.PlayerData {
	pd := structs.PlayerData{
		Damage: structs.PlayerDamage{
			Crit:       round(float64(entity.SkillStats.Crits) / float64(entity.SkillStats.Hits) * 100),
			CritDamage: round(float64(entity.DamageStats.CritDamage) / float64(entity.DamageStats.Damage) * 100),
			FA:         round(float64(entity.DamageStats.FADamage) / float64(entity.DamageStats.Damage) * 100),
			BA:         round(float64(entity.DamageStats.BADamage) / float64(entity.DamageStats.Damage) * 100),
			Buff:       round(float64(entity.DamageStats.Buffed) / float64(entity.DamageStats.Damage) * 100),
			Brand:      round(float64(entity.DamageStats.Debuffed) / float64(entity.DamageStats.Damage) * 100),
		},
	}

	buffs := Buffs{}
	buffs.Collect(
		enc.raw.DamageStats.Buffs,
		entity.DamageStats.BuffedBy,
		DefaultBuffFilter,
	)
	buffs.Collect(
		enc.raw.DamageStats.Debuffs,
		entity.DamageStats.DebuffedBy,
		DefaultBuffFilter,
	)
	for _, buff := range buffs {
		buff.Percent = round(float64(buff.Damage) / float64(entity.DamageStats.Damage) * 100)
	}
	pd.Synergy = structs.Buffs(buffs)

	skillDamage := make(map[string]structs.SkillDamage)
	skillBuffs := make(map[string]structs.Buffs)
	for id, skill := range entity.Skills {
		skillDamage[id] = Skill(enc.raw, entity, skill)
		if skill.Damage == 0 {
			continue
		}

		buffs := Buffs{}
		buffs.Collect(
			enc.raw.DamageStats.Buffs,
			skill.BuffedBy,
			DefaultBuffFilter,
		)
		buffs.Collect(
			enc.raw.DamageStats.Debuffs,
			skill.DebuffedBy,
			DefaultBuffFilter,
		)
		for _, buff := range buffs {
			buff.Percent = round(float64(buff.Damage) / float64(skill.Damage) * 100)
		}
		skillBuffs[id] = structs.Buffs(buffs)
	}
	pd.SkillDamage = skillDamage
	pd.SkillSynergy = skillBuffs

	return pd
}

func Skill(enc *meter.Encounter, player meter.Entity, skill meter.Skill) structs.SkillDamage {
	return structs.SkillDamage{
		Damage:     skill.Damage,
		DPS:        skill.DPS,
		Percent:    round(float64(skill.Damage) / float64(player.DamageStats.Damage) * 100),
		Crit:       round(float64(skill.Crits) / float64(skill.Hits) * 100),
		CritDamage: round(float64(skill.CritDamage) / float64(skill.Damage) * 100),
		FA:         round(float64(skill.FADamage) / float64(skill.Damage) * 100),
		BA:         round(float64(skill.BADamage) / float64(skill.Damage) * 100),
		Buff:       round(float64(skill.Buffed) / float64(skill.Damage) * 100),
		Brand:      round(float64(skill.Debuffed) / float64(skill.Damage) * 100),
		APH:        round(float64(skill.Damage) / float64(skill.Hits) * 100),
		APC:        round(float64(skill.Damage) / float64(skill.Casts) * 100),
		Max:        skill.Max,
		Casts:      skill.Casts,
		CPM:        round(float64(skill.Casts) / (float64(enc.Duration) / 1000 / 60)),
		Hits:       skill.Hits,
		HPM:        round(float64(skill.Hits) / (float64(enc.Duration) / 1000 / 60)),
	}
}

type BuffFilter func(info meter.Buff) bool

func DefaultBuffFilter(info meter.Buff) bool {
	return slices.Contains(
		[]string{"classskill", "identity", "ability"}, info.BuffCategory,
	) && info.Target == "PARTY" && (1|2|4|128)&info.BuffType != 0
}

type Buffs map[string]*structs.BuffGroup

func (b Buffs) Collect(catalog meter.BuffInfo, buffs meter.BuffDamage, filter BuffFilter) {
	for buff, damage := range buffs {
		info, ok := catalog[buff]
		if !ok || !filter(info) {
			continue
		}

		group := "0_"
		if info.Source.Skill != nil {
			group = fmt.Sprintf("%d_", info.Source.Skill.ClassID)
		}

		if info.UniqueGroup != 0 {
			group += strconv.Itoa(int(info.UniqueGroup))
		} else if info.Source.Skill != nil {
			group += info.Source.Skill.Name
		} else {
			// uh oh
		}

		entry, ok := b[group]
		if !ok {
			entry = &structs.BuffGroup{
				Buffs: map[string]int64{},
			}
		}
		entry.Buffs[buff] = damage
		entry.Damage += damage

		b[group] = entry
	}
}

func (enc *Encounter) highlight() {

}

func (p *Processor) Save(ctx context.Context, user pgtype.UUID, str string, raw *meter.Encounter) (int32, error) {
	tx, err := p.db.Pool.Begin(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "begin transaction")
	}
	defer tx.Rollback(ctx)
	qtx := p.db.Queries.WithTx(tx)

	enc, err := p.Process(raw)
	if err != nil {
		return 0, errors.Wrap(err, "processing encounter")
	}

	start := time.UnixMilli(raw.FightStart).UTC()
	var date pgtype.Timestamp
	if err := date.Scan(start); err != nil {
		return 0, errors.Wrap(err, "scanning duration pgtype.Timstamp")
	}

	encID, err := qtx.InsertEncounter(ctx, sql.InsertEncounterParams{
		UploadedBy:  user,
		Boss:        raw.CurrentBossName,
		Difficulty:  raw.Difficulty,
		Date:        date,
		Duration:    raw.Duration,
		LocalPlayer: raw.LocalPlayer,
		Header:      enc.Header,
		Data:        enc.Data,
	})
	if err != nil {
		return 0, errors.Wrap(err, "inserting encounter")
	}

	if err := tx.Commit(ctx); err != nil {
		return 0, errors.Wrap(err, "committing transaction")
	}

	if err := p.s3.SaveEncounter(ctx, encID, str); err != nil {
		return 0, errors.Wrap(err, "saving encounter to s3")
	}

	return encID, nil
}
