package structs

import "github.com/0fau/logs/pkg/process/meter"

type PlayerHeader struct {
	Name      string    `json:"name"`
	Class     string    `json:"class"`
	GearScore float64   `json:"gearScore"`
	Damage    int64     `json:"damage"`
	Percent   string    `json:"percent"`
	DPS       int64     `json:"dps"`
	Dead      bool      `json:"dead"`
	DeadFor   int64     `json:"deadFor"`
	Highlight []float64 `json:"highlight"`
}

type EncounterHeader struct {
	Players map[string]PlayerHeader `json:"players"`
	Parties [][]string              `json:"parties"`
	Damage  int64                   `json:"damage"`
	Cleared bool                    `json:"cleared"`
}

type PartyDamageRow []float64

type EncounterData struct {
	BuffCatalog  map[string]BuffInfo  `json:"buffCatalog"`
	SkillCatalog map[string]SkillInfo `json:"skillCatalog"`

	Synergies [][]BuffGroupInfo     `json:"synergies"`
	SelfBuffs []BuffGroupInfo       `json:"selfBuffs"`
	Players   map[string]PlayerData `json:"players"`

	BossHPLog meter.HPLog `json:"bossHPLog"`
}

type BuffGroupInfo struct {
	Name  string   `json:"name"`
	Buffs []string `json:"buffs"`
}

type BuffInfo struct {
	Name        string     `json:"name"`
	Icon        string     `json:"icon"`
	Description string     `json:"description"`
	Category    string     `json:"category"`
	Set         string     `json:"set"`
	Skill       *BuffSkill `json:"skill"`
}

type BuffSkill struct {
	Class       int    `json:"class"`
	Description string `json:"description"`
	Name        string `json:"name"`
	Icon        string `json:"icon"`
	ID          int    `json:"id"`
}

type SkillInfo struct {
	Name string `json:"name"`
	Icon string `json:"icon"`
}

type Buffs map[string]*BuffGroup

type PlayerData struct {
	Damage      PlayerDamage           `json:"damage"`
	SkillDamage map[string]SkillDamage `json:"skillDamage"`

	Synergy      Buffs            `json:"synergy"`
	SkillSynergy map[string]Buffs `json:"skillSynergy"`

	SelfBuff       Buffs            `json:"selfBuff"`
	SkillSelfBuffs []BuffGroupInfo  `json:"skillSelfBuffs"`
	SkillSelfBuff  map[string]Buffs `json:"skillSelfBuff"`

	DPSLog []int64 `json:"dpsLog"`
}

type PlayerDamage struct {
	Crit       string `json:"crit"`
	CritDamage string `json:"critDamage"`
	FA         string `json:"fa"`
	BA         string `json:"ba"`
	Buff       string `json:"buff"`
	Brand      string `json:"brand"`
	Identity   string `json:"identity"`
	Casts      int32  `json:"casts"`
	CPM        string `json:"cpm"`
	Hits       int32  `json:"hits"`
	HPM        string `json:"hpm"`
}

type SkillDamage struct {
	Damage     int64  `json:"damage"`
	DPS        int64  `json:"dps"`
	Percent    string `json:"percent"`
	Crit       string `json:"crit"`
	CritDamage string `json:"critDamage"`
	FA         string `json:"fa"`
	BA         string `json:"ba"`
	Buff       string `json:"buff"`
	Brand      string `json:"brand"`
	APH        string `json:"aph"`
	APC        string `json:"apc"`
	Max        int64  `json:"max"`
	Casts      int32  `json:"casts"`
	CPM        string `json:"cpm"`
	Hits       int32  `json:"hits"`
	HPM        string `json:"hpm"`

	CastLog     []int32           `json:"castLog"`
	TripodIndex *meter.TripodRows `json:"tripodIndex"`
	TripodLevel *meter.TripodRows `json:"tripodLevel"`
}

type BuffGroup struct {
	Damage  int64            `json:"damage"`
	Percent string           `json:"percent"`
	Buffs   map[string]int64 `json:"buffs"`
}

type EncounterSettings struct {
	Visibility []string `json:"visibility"`
}
