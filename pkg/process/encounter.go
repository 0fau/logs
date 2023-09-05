package process

type RawEncounter struct {
	CurrentBossName string               `json:"currentBossName"`
	Duration        int32                `json:"duration"`
	Entities        map[string]RawEntity `json:"entities"`
	FightStart      int64                `json:"fightStart"`
	DamageStats     EncounterDamageStats `json:"encounterDamageStats"`
}

type EncounterDamageStats struct {
	TotalDamageDealt int64                    `json:"totalDamageDealt"`
	Misc             EncounterDamageStatsMisc `json:"misc"`
}

type EncounterDamageStatsMisc struct {
	Cleared bool `json:"raidClear"`
}

type RawEntity struct {
	ClassId     int               `json:"classId"`
	Class       string            `json:"class"`
	Name        string            `json:"name"`
	EntityType  string            `json:"entityType"`
	DamageStats EntityDamageStats `json:"damageStats"`
}

type EntityDamageStats struct {
	BackAttackDamage int   `json:"backAttackDamage"`
	DamageDealt      int64 `json:"damageDealt"`
	DPS              int32 `json:"dps"`
}
