package process

import (
	"fmt"
	"os"

	"github.com/cockroachdb/errors"
	"github.com/goccy/go-json"
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

func (p *Processor) loadMeterData() error {
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
