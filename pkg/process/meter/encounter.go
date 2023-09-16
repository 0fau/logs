package meter

type Encounter struct {
	CurrentBossName string               `json:"currentBossName"`
	Duration        int32                `json:"duration"`
	Entities        map[string]Entity    `json:"entities"`
	FightStart      int64                `json:"fightStart"`
	DamageStats     EncounterDamageStats `json:"encounterDamageStats"`
	LocalPlayer     string               `json:"localPlayer"`
}

type EncounterDamageStats struct {
	TotalDamageDealt int64                    `json:"totalDamageDealt"`
	Buffs            map[string]Buff          `json:"buffs"`
	Debuffs          map[string]Buff          `json:"debuffs"`
	Misc             EncounterDamageStatsMisc `json:"misc"`
}

type EncounterDamageStatsMisc struct {
	Cleared bool `json:"raidClear"`
}

type Entity struct {
	ClassId     int32             `json:"classId"`
	Class       string            `json:"class"`
	Name        string            `json:"name"`
	EntityType  string            `json:"entityType"`
	DamageStats EntityDamageStats `json:"damageStats"`
	Skills      map[string]Skill  `json:"skills"`
}

type EntityDamageStats struct {
	BackAttackDamage int32            `json:"backAttackDamage"`
	DamageDealt      int64            `json:"damageDealt"`
	DPS              int64            `json:"dps"`
	BuffedBy         map[string]int64 `json:"buffedBy"`
	DebuffedBy       map[string]int64 `json:"debuffedBy"`
}
