package meter

type Buff struct {
	BuffCategory string `json:"buffCategory"`
	BuffType     int32  `json:"buffType"`
	Category     string `json:"category"`
	Target       string `json:"target"`
	UniqueGroup  int32  `json:"uniqueGroup"`

	Source BuffSource `json:"source"`
}

type BuffSource struct {
	Description       string     `json:"desc"`
	Icon              string     `json:"string"`
	Name              string     `json:"name"`
	SetName           *string    `json:"setName"`
	Skill             *BuffSkill `json:"skill"`
	SummonIDs         []int32    `json:"summonIds"`
	SummonSourceSkill []int32    `json:"summonSourceSkill"`
	SourceSkill       int32      `json:"sourceSkill"`
}

type BuffSkill struct {
	ClassID     int32  `json:"classId"`
	Description string `json:"desc"`
	Icon        string `json:"icon"`
	ID          int32  `json:"id"`
	Name        string `json:"name"`
}

type Skill struct {
	BuffedBy    map[string]int64 `json:"buffedBy"`
	CastLog     []int32          `json:"castLog"`
	Casts       int32            `json:"casts"`
	Crits       int32            `json:"crits"`
	DebuffedBy  map[string]int64 `json:"debuffedBy"`
	DPS         int64            `json:"dps"`
	Hits        int32            `json:"hits"`
	Icon        string           `json:"icon"`
	MaxDamage   int64            `json:"maxDamage"`
	TotalDamage int64            `json:"totalDamage"`
	Name        string           `json:"name"`
}
