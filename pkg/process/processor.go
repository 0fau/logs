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
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"hash/fnv"
	"math"
	"slices"
	"strconv"
	"strings"
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
	//if err := p.loadMeterData(); err != nil {
	//	return err
	//}

	return nil
}

type Encounter struct {
	Raw *meter.Encounter

	Header structs.EncounterHeader
	Data   structs.EncounterData
}

func (p *Processor) Process(raw *meter.Encounter) (*Encounter, error) {
	enc := &Encounter{Raw: raw}
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
		Parties: make([][]string, len(enc.Raw.DamageStats.Misc.PartyInfo)),
		Damage:  enc.Raw.DamageStats.TotalDamageDealt,
		Cleared: enc.Raw.DamageStats.Misc.Cleared,
	}

	for party, players := range enc.Raw.DamageStats.Misc.PartyInfo {
		num, err := strconv.Atoi(party)
		if err != nil {
			return structs.EncounterHeader{}, errors.Wrapf(err, "converting party %s to number", party)
		}

		slices.SortFunc(players, func(a, b string) int {
			return cmp.Compare(
				enc.Raw.Entities[b].DamageStats.Damage,
				enc.Raw.Entities[a].DamageStats.Damage,
			)
		})
		header.Parties[num] = players
	}

	for _, entity := range enc.Raw.Entities {
		if entity.EntityType != "PLAYER" {
			continue
		}

		header.Players[entity.Name] = structs.PlayerHeader{
			Name:      entity.Name,
			Class:     entity.Class,
			GearScore: entity.GearScore,
			Damage:    entity.DamageStats.Damage,
			Percent:   round(float64(entity.DamageStats.Damage) / float64(enc.Raw.DamageStats.TotalDamageDealt) * 100),
			DPS:       entity.DamageStats.DPS,
			Dead:      entity.Dead,
			DeadFor:   enc.Raw.End - entity.DamageStats.DeathTime,
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

type BuffGroups map[string]BuffGroup

type BuffGroup struct {
	Name  string
	Buffs map[string]struct{}
}

func (bgs BuffGroups) Collect(group, buff string) {
	bg, ok := bgs[group]
	if !ok {
		bg.Name = group
		bg.Buffs = make(map[string]struct{})
	}
	bg.Buffs[buff] = struct{}{}
	bgs[group] = bg
}

func (bgs BuffGroups) Serialize(order func(a, b structs.BuffGroupInfo) int) []structs.BuffGroupInfo {
	arr := make([]structs.BuffGroupInfo, 0, len(bgs))
	for _, bg := range bgs {
		ssyn := structs.BuffGroupInfo{
			Name:  bg.Name,
			Buffs: make([]string, 0, len(bg.Buffs)),
		}
		for buff := range bg.Buffs {
			ssyn.Buffs = append(ssyn.Buffs, buff)
		}
		slices.SortFunc(ssyn.Buffs, func(a, b string) int {
			num1, _ := strconv.Atoi(a)
			num2, _ := strconv.Atoi(b)
			return cmp.Compare(num1, num2)
		})
		arr = append(arr, ssyn)
	}
	slices.SortFunc(arr, order)
	return arr
}

var SynergyOrder = map[string]int{
	"204_101204":              0, // Bard ATK
	"204_210230":              1, // Bard Brand
	"204_Serenade of Courage": 2,
	"204_Heavenly Tune":       3,

	"105_101105": 4, // Pally ATK
	"105_210230": 5, // Pally Brand
	"105_368000": 6, // Pally Aura

	"602_314004": 6, // Artist ATK
	"602_210230": 7, // Artist Brand
	"602_310501": 8, // Artist Moonfall
}

func compareSynergy(a, b structs.BuffGroupInfo) int {
	aorder, aok := SynergyOrder[a.Name]
	border, bok := SynergyOrder[b.Name]
	if aok && !bok {
		return -1
	} else if !aok && bok {
		return 1
	} else if aok && bok {
		return cmp.Compare(aorder, border)
	}

	asplit := strings.Split(a.Name, "_")
	bsplit := strings.Split(b.Name, "_")
	if asplit[0] != bsplit[0] {
		return cmp.Compare(asplit[0], bsplit[0])
	}

	return cmp.Compare(asplit[1], bsplit[1])
}

func compareBuff(a, b structs.BuffGroupInfo) int {
	return cmp.Compare(a.Name, b.Name)
}

func (enc *Encounter) processData() (structs.EncounterData, error) {
	data := structs.EncounterData{
		Players:      make(map[string]structs.PlayerData),
		Synergies:    make([][]structs.BuffGroupInfo, len(enc.Header.Parties)),
		BuffCatalog:  make(map[string]structs.BuffInfo),
		SkillCatalog: make(map[string]structs.SkillInfo),
		BossHPLog:    enc.Raw.DamageStats.Misc.HPLog,
	}

	parties := map[string]int{}
	for party, players := range enc.Header.Parties {
		for _, player := range players {
			parties[player] = party
		}
	}

	partyBuffs := make([]BuffGroups, len(enc.Header.Parties))
	for i := range partyBuffs {
		partyBuffs[i] = make(BuffGroups)
	}
	selfBuffs := make(BuffGroups)

	for name, entity := range enc.Raw.Entities {
		if entity.EntityType != "PLAYER" {
			continue
		}

		data.Players[name] = enc.processPlayer(entity)
		for gname, group := range data.Players[name].Synergy {
			for buff := range group.Buffs {
				partyBuffs[parties[name]].Collect(gname, buff)
				enc.CatalogBuff(data, buff)
			}
		}

		for gname, group := range data.Players[name].SelfBuff {
			for buff := range group.Buffs {
				selfBuffs.Collect(gname, buff)
				enc.CatalogBuff(data, buff)
			}
		}

		skillSelfBuffs := make(BuffGroups)
		for _, groups := range data.Players[name].SkillSelfBuff {
			for gname, group := range groups {
				for buff := range group.Buffs {
					enc.CatalogBuff(data, buff)
					skillSelfBuffs.Collect(gname, buff)
				}
			}
		}
		player, ok := data.Players[name]
		if ok {
			player.SkillSelfBuffs = skillSelfBuffs.Serialize(compareBuff)
			data.Players[name] = player
		}

		for skill := range data.Players[name].SkillDamage {
			info := enc.Raw.Entities[name].Skills[skill]
			data.SkillCatalog[skill] = structs.SkillInfo{
				Name: info.Name,
				Icon: info.Icon,
			}
		}
	}

	for i, groups := range partyBuffs {
		data.Synergies[i] = groups.Serialize(compareSynergy)
	}
	data.SelfBuffs = selfBuffs.Serialize(compareBuff)

	return data, nil
}

type Player struct {
}

func (enc *Encounter) ProcessPlayer() Player {
	return Player{}
}

func (enc *Encounter) CatalogBuff(data structs.EncounterData, buff string) {
	if _, ok := data.BuffCatalog[buff]; ok {
		return
	}

	info, err := enc.BuffInfo(buff)
	if err != nil {
		return
	}

	data.BuffCatalog[buff] = info
}

func (enc *Encounter) BuffInfo(buff string) (structs.BuffInfo, error) {
	info, ok := enc.Raw.DamageStats.Buffs[buff]
	if !ok {
		info, ok = enc.Raw.DamageStats.Debuffs[buff]
		if !ok {
			return structs.BuffInfo{}, errors.New("buff info not found")
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
	return binfo, nil
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
			Casts:      entity.SkillStats.Casts,
			CPM:        round(float64(entity.SkillStats.Casts) / (float64(enc.Raw.Duration) / 1000 / 60)),
			Hits:       entity.SkillStats.Hits,
			HPM:        round(float64(entity.SkillStats.Hits) / (float64(enc.Raw.Duration) / 1000 / 60)),
		},
		DPSLog: entity.DamageStats.DPSAverage,
	}

	catalogs := []meter.BuffInfo{
		enc.Raw.DamageStats.Buffs,
		enc.Raw.DamageStats.Debuffs,
	}

	buffs, self, skillSelf := Buffs{}, Buffs{}, Buffs{}
	for i, damage := range []meter.BuffDamage{
		entity.DamageStats.BuffedBy,
		entity.DamageStats.DebuffedBy,
	} {
		catalog := catalogs[i]
		buffs.CollectAll(
			catalog, damage,
			PartySynergyFilter,
		)
		self.CollectAll(
			catalog, damage,
			SelfBuffFilter,
		)
		skillSelf.CollectAll(
			catalog, damage,
			PlayerSelfBuffFilter(entity),
		)
	}
	for _, buffs := range []Buffs{buffs, self, skillSelf} {
		for _, buff := range buffs {
			buff.Percent = round(float64(buff.Damage) / float64(entity.DamageStats.Damage) * 100)
		}
	}
	pd.Synergy = structs.Buffs(buffs)
	pd.SelfBuff = structs.Buffs(self)

	skillDamage := make(map[string]structs.SkillDamage)
	skillBuffs := make(map[string]structs.Buffs)
	skillSelfBuffs := make(map[string]structs.Buffs)
	for id, skill := range entity.Skills {
		skillDamage[id] = Skill(enc.Raw, entity, skill)
		if skill.Damage == 0 {
			continue
		}

		buffs, self := Buffs{}, Buffs{}
		for i, damage := range []meter.BuffDamage{
			skill.BuffedBy,
			skill.DebuffedBy,
		} {
			catalog := catalogs[i]
			buffs.CollectAll(
				catalog, damage,
				PartySynergyFilter,
			)
			self.CollectAll(
				catalog, damage,
				PlayerSelfBuffFilter(entity),
			)
		}
		for _, buff := range buffs {
			buff.Percent = round(float64(buff.Damage) / float64(skill.Damage) * 100)
		}
		for _, buff := range self {
			buff.Percent = round(float64(buff.Damage) / float64(skill.Damage) * 100)
		}
		skillBuffs[id] = structs.Buffs(buffs)
		skillSelfBuffs[id] = structs.Buffs(self)
	}
	pd.SkillDamage = skillDamage
	pd.SkillSynergy = skillBuffs

	skillSelfBuffs["_player"] = structs.Buffs(skillSelf)
	pd.SkillSelfBuff = skillSelfBuffs

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

type BuffFilter func(info meter.Buff) (string, bool)

func PartySynergyFilter(info meter.Buff) (string, bool) {
	if !(slices.Contains(
		[]string{"classskill", "identity", "ability"}, info.BuffCategory,
	) && info.Target == "PARTY" && (1|2|4|128)&info.BuffType != 0) {
		return "", false
	}

	return BuffGroupName(info), true
}

func BuffGroupName(info meter.Buff) string {
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

	return group
}

func SelfBuffFilter(info meter.Buff) (string, bool) {
	if info.Target == "PARTY" || (1|2|4|128)&info.BuffType == 0 {
		return "", false
	}

	var group string
	switch info.BuffCategory {
	case "set":
		group = "set_" + info.Source.SetName
	case "bracelet", "elixir":
		group = fmt.Sprintf("%s_%d", info.BuffCategory, info.UniqueGroup)
	case "pet", "cook", "battleitem", "dropsofether":
		group = info.BuffCategory
	default:
		return "", false
	}

	return group, true
}

func PlayerSelfBuffFilter(player meter.Entity) BuffFilter {
	return func(info meter.Buff) (string, bool) {
		if info.Target == "PARTY" || (1|2|4|128)&info.BuffType == 0 {
			return "", false
		}

		var group string
		switch info.BuffCategory {
		case "ability":
			if info.UniqueGroup != 0 {
				group = fmt.Sprintf("%d", info.UniqueGroup)
			}
		case "etc":
			group = "etc_" + info.Source.Name
		case "classskill", "identity":
			if info.Source.Skill != nil &&
				player.ClassId != info.Source.Skill.ClassID {
				return "", false
			}

			group = BuffGroupName(info)
		default:
			return "", false
		}

		return group, true
	}
}

type Buffs map[string]*structs.BuffGroup

func (b Buffs) CollectAll(catalog meter.BuffInfo, damages meter.BuffDamage, filter BuffFilter) {
	for buff, damage := range damages {
		info, ok := catalog[buff]
		if !ok {
			continue
		}

		group, ok := filter(info)
		if !ok {
			continue
		}
		if group == "" {
			group = buff
		}

		b.Collect(group, buff, damage)
	}
}

func (b Buffs) Collect(group, buff string, damage int64) {
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

func (enc *Encounter) highlight() {

}

func (enc *Encounter) UniqueHash(players []string) string {
	slices.Sort(players)

	var buf strings.Builder
	buf.WriteString(enc.Raw.CurrentBossName)
	buf.WriteString(" ")
	buf.WriteString(enc.Raw.Difficulty)

	for _, player := range players {
		buf.WriteString(" ")
		buf.WriteString(player)
	}

	h := fnv.New32a()
	h.Write([]byte(buf.String()))
	return fmt.Sprintf("%d", h.Sum32())
}

func (p *Processor) Save(ctx context.Context, user pgtype.UUID, str string, raw *meter.Encounter) (int32, error) {
	enc, err := p.Process(raw)
	if err != nil {
		return 0, errors.Wrap(err, "processing encounter")
	}

	order := make([]string, 0, len(enc.Header.Players))
	for player := range enc.Header.Players {
		order = append(order, player)
	}
	slices.SortFunc(order, func(a, b string) int {
		return cmp.Compare(
			enc.Header.Players[b].Damage,
			enc.Header.Players[a].Damage,
		)
	})
	hash := enc.UniqueHash(order)

	start := time.UnixMilli(raw.FightStart).UTC()
	var date pgtype.Timestamp
	if err := date.Scan(start); err != nil {
		return 0, errors.Wrap(err, "scanning duration pgtype.Timstamp")
	}

	var encID int32
	if err := crdbpgx.ExecuteTx(ctx, p.db.Pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		qtx := p.db.Queries.WithTx(tx)

		group, err := qtx.GetUniqueGroup(ctx, sql.GetUniqueGroupParams{
			UniqueHash: hash,
			Date:       date,
		})
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return errors.Wrap(err, "getting unique group")
		}

		encID, err = qtx.InsertEncounter(ctx, sql.InsertEncounterParams{
			UploadedBy:  user,
			Boss:        raw.CurrentBossName,
			Difficulty:  raw.Difficulty,
			Date:        date,
			Duration:    raw.Duration,
			LocalPlayer: raw.LocalPlayer,
			Header:      enc.Header,
			Data:        enc.Data,
			Version:     1,
			UniqueHash:  hash,
			UniqueGroup: group,
		})
		if err != nil {
			return errors.Wrap(err, "inserting encounter")
		}

		if group == 0 {
			if err := qtx.UpdateUniqueGroup(ctx, sql.UpdateUniqueGroupParams{
				ID:          encID,
				UniqueGroup: encID,
			}); err != nil {
				return errors.Wrap(err, "updating unique group")
			}
		} else if group != encID {
			if err := qtx.UpsertEncounterGroup(ctx, sql.UpsertEncounterGroupParams{
				GroupID: group,
				Column2: user,
			}); err != nil {
				return errors.Wrap(err, "upserting encounter group")
			}
		}

		if err := qtx.InsertCharacter(ctx, sql.InsertCharacterParams{
			UserID:    user,
			Character: enc.Raw.LocalPlayer,
			Class:     enc.Header.Players[raw.LocalPlayer].Class,
			GearScore: enc.Header.Players[raw.LocalPlayer].GearScore,
		}); err != nil {
			return errors.Wrap(err, "inserting character")
		}

		players := make([]sql.InsertPlayerParams, 0, len(enc.Header.Players))
		for i, player := range order {
			header := enc.Header.Players[player]
			players = append(players, sql.InsertPlayerParams{
				Encounter:  encID,
				Boss:       raw.CurrentBossName,
				Difficulty: raw.Difficulty,
				Class:      header.Class,
				Name:       player,
				Dead:       header.Dead,
				Place:      int32(i + 1),
				GearScore:  header.GearScore,
				Dps:        header.DPS,
			})
		}
		if _, err := qtx.InsertPlayer(ctx, players); err != nil {
			return errors.Wrap(err, "inserting players")
		}

		return nil
	}); err != nil {
		return 0, errors.Wrap(err, "executing transaction")
	}

	if err := p.s3.SaveEncounter(ctx, encID, str); err != nil {
		return 0, errors.Wrap(err, "saving encounter to s3")
	}

	return encID, nil
}
