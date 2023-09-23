package process

import (
	"encoding/json"
	"fmt"
	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/process/meter"
	"github.com/cockroachdb/errors"
	"os"
)

type SkillData struct {
	ID                int32  `json:"id"`
	Name              string `json:"name"`
	Desc              string `json:"desc"`
	ClassID           int32  `json:"classid"`
	Icon              string `json:"icon"`
	SummonSourceSkill int32  `json:"summonsourceskill"`
	SourceSkill       int32  `json:"sourceskill"`
}

type BuffData struct {
	ID            int32           `json:"id"`
	Name          string          `json:"name"`
	Desc          string          `json:"desc"`
	Icon          string          `json:"icon"`
	IconShowType  string          `json:"iconshowtype"`
	Duration      int32           `json:"duration"`
	BuffType      string          `json:"type"`
	BuffCategory  string          `json:"buffcategory"`
	Target        string          `json:"target"`
	UniqueGroup   int32           `json:"uniquegroup"`
	OverlapFlag   int32           `json:"overlapflag"`
	PassiveOption []PassiveOption `json:"passiveoption"`
	SourceSkill   *int32          `json:"sourceskill"`
	SetName       *string         `json:"setname"`
}

type PassiveOption struct {
	Type     string `json:"type"`
	KeyStat  string `json:"keystat"`
	KeyIndex int32  `json:"keyindex"`
	Value    int32  `json:"value"`
}

type LogProcessor struct {
	skills map[string]SkillData
	buffs  map[string]BuffData
}

func NewLogProcessor() *LogProcessor {
	return &LogProcessor{}
}

func (p *LogProcessor) readMeterData() error {
	data := []struct {
		file string
		dest interface{}
	}{
		{"Skill", &p.skills},
		{"SkillBuff", &p.buffs},
	}

	for _, entry := range data {
		raw, err := os.ReadFile(fmt.Sprintf("meter-data/%s.json", entry.file))
		if err != nil {
			return errors.Wrapf(err, "reading meter data: %s", entry.file)
		}
		if err := json.Unmarshal(raw, entry.dest); err != nil {
			return errors.Wrapf(err, "unmarshalling meter data: %s", entry.file)
		}
	}

	return nil
}

func (p *LogProcessor) Initialize() error {
	if err := p.readMeterData(); err != nil {
		return err
	}

	return nil
}

type Encounter struct {
}

func (p *LogProcessor) Lint(raw *meter.Encounter) error {
	return nil
}

func (p *LogProcessor) Fetch(db *database.DB, encounter int32) (*Encounter, error) {
	return nil, nil
}

func (p *LogProcessor) Process() {

}
