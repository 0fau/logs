package meter

type BuffInfo map[string]Buff

type Encounter struct {
	CurrentBossName string               `json:"currentBossName"`
	Difficulty      string               `json:"difficulty"`
	Duration        int32                `json:"duration"`
	Entities        map[string]Entity    `json:"entities"`
	FightStart      int64                `json:"fightStart"`
	End             int64                `json:"lastCombatPacket"`
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

type Engraving struct {
	Name  string `json:"name"`
	ID    int32  `json:"id"`
	Level int32  `json:"level"`
	Icon  string `json:"icon"`
}

type EngravingData struct {
	ClassEngravings []Engraving `json:"classEngravings"`
	OtherEngravings []Engraving `json:"otherEngravings"`
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
	Engravings  *EngravingData    `json:"engravingData"`
	Skills      map[string]Skill  `json:"skills"`
}

type EntitySkillStats struct {
	Hits  int32 `json:"hits"`
	Casts int32 `json:"casts"`
	BA    int32 `json:"backAttacks"`
	FA    int32 `json:"frontAttacks"`
	Crits int32 `json:"crits"`

	IdentityStats string `json:"identityStats"`
}

type BuffDamage map[string]int64

type EntityDamageStats struct {
	Damage     int64   `json:"damageDealt"`
	DPS        int64   `json:"dps"`
	DPSAverage []int64 `json:"dpsAverage"`
	DPSRolling []int64 `json:"dpsRolling10sAvg"`
	CritDamage int64   `json:"critDamage"`
	BADamage   int64   `json:"backAttackDamage"`
	FADamage   int64   `json:"frontAttackDamage"`

	RDPSDamageReceived        int64 `json:"rdpsDamageReceived"`
	RDPSDamageReceivedSupport int64 `json:"rdpsDamageReceivedSupport"`
	RDPSDamageGiven           int64 `json:"rdpsDamageGiven"`

	Buffed         int64      `json:"buffedBySupport"`
	BuffedIdentity int64      `json:"buffedByIdentity"`
	Debuffed       int64      `json:"debuffedBySupport"`
	BuffedBy       BuffDamage `json:"buffedBy"`
	DebuffedBy     BuffDamage `json:"debuffedBy"`
	DeathTime      int64      `json:"deathTime"`

	ShieldsGivenBy           BuffDamage `json:"shieldsGivenBy"`
	ShieldsReceivedBy        BuffDamage `json:"ShieldsReceivedBy"`
	DamageAbsorbedBy         BuffDamage `json:"damageAbsorbedBy"`
	DamageAbsorbedOnOthersBy BuffDamage `json:"damageAbsorbedOnOthersBy"`
}
