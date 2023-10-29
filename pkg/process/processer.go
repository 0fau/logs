package process

import (
	"context"
	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/0fau/logs/pkg/process/meter"
	"github.com/cockroachdb/errors"
	"github.com/jackc/pgx/v5/pgtype"
	"strconv"
	"time"
)

type Processor struct {
	skills map[string]SkillData
	buffs  map[string]BuffData
}

func NewLogProcessor() *Processor {
	return &Processor{}
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
	data, err := enc.processData()
	if err != nil {
		return nil, errors.Wrap(err, "processing encounter data")
	}

	enc.Header = header
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

const (
	PartyDamagePercentColumn = iota
	PartyDamageCritColumn
	PartyDamageCritDamageColumn
	PartyDamageFAColumn
	PartyDamageBAColumn
	PartyDamageBuffPercentColumn
	PartyDamageBrandPercentColumn
)

func (enc *Encounter) processData() (structs.EncounterData, error) {
	data := structs.EncounterData{}

	for _, entity := range enc.raw.Entities {

	}

	return data, nil
}

func (enc *Encounter) highlight() {

}

func (p *Processor) Save(ctx context.Context, db *database.DB, user pgtype.UUID, raw *meter.Encounter) (int32, error) {
	tx, err := db.Pool.Begin(ctx)
	if err != nil {
		return 0, errors.Wrap(err, "begin transaction")
	}
	defer tx.Rollback(ctx)
	qtx := db.Queries.WithTx(tx)

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
	return encID, nil
}
