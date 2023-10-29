package meter

type BuffInfo map[string]Buff

type Encounter struct {
	CurrentBossName string               `json:"currentBossName"`
	Difficulty      string               `json:"difficulty"`
	Duration        int32                `json:"duration"`
	Entities        map[string]Entity    `json:"entities"`
	FightStart      int64                `json:"fightStart"`
	DamageStats     EncounterDamageStats `json:"encounterDamageStats"`
	LocalPlayer     string               `json:"localPlayer"`
}

type EncounterDamageStats struct {
	TotalDamageDealt int64                    `json:"totalDamageDealt"`
	Buffs            BuffInfo                 `json:"buffs"`
	Debuffs          BuffInfo                 `json:"debuffs"`
	Misc             EncounterDamageStatsMisc `json:"misc"`
}

type HPLogEntry struct {
	HP   int64   `json:"hp"`
	P    float64 `json:"p"`
	Time int32   `json:"time"`
}

type PartyInfo map[string][]string

type HPLog map[string][]HPLogEntry

type EncounterDamageStatsMisc struct {
	Cleared   bool      `json:"raidClear"`
	HPLog     HPLog     `json:"bossHpLog"`
	PartyInfo PartyInfo `json:"partyInfo"`
}

type Entity struct {
	ClassId     int32             `json:"classId"`
	Class       string            `json:"class"`
	Dead        bool              `json:"isDead"`
	Name        string            `json:"name"`
	GearScore   float64           `json:"gearScore"`
	EntityType  string            `json:"entityType"`
	DamageStats EntityDamageStats `json:"damageStats"`
	SkillStats  EntitySkillStats  `json:"skillStats"`
	Skills      map[string]Skill  `json:"skills"`
}

type EntitySkillStats struct {
	Hits  int32 `json:"hits"`
	Casts int32 `json:"casts"`
	BA    int32 `json:"backAttacks"`
	FA    int32 `json:"frontAttacks"`
	Crits int32 `json:"crits"`
}

type BuffDamage map[string]int64

type EntityDamageStats struct {
	Damage     int64      `json:"damageDealt"`
	DPS        int64      `json:"dps"`
	DPSAverage []int64    `json:"dpsAverage"`
	DPSRolling []int64    `json:"dpsRolling10sAvg"`
	BADamage   int64      `json:"backAttackDamage"`
	FADamage   int64      `json:"frontAttackDamage"`
	Buffed     int64      `json:"buffedBySupport"`
	Debuffed   int64      `json:"debuffedBySupport"`
	BuffedBy   BuffDamage `json:"buffedBy"`
	DebuffedBy BuffDamage `json:"debuffedBy"`
	DeathTime  int64      `json:"deathTime"`
}
