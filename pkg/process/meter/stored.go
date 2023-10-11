package meter

type StoredEncounterFields struct {
	Buffs     BuffInfo                `json:"buffs"`
	Debuffs   BuffInfo                `json:"debuffs"`
	PartyInfo map[string][]string     `json:"partyInfo"`
	HPLog     map[string][]HPLogEntry `json:"bossHpLog"`
}

type StoredEntityFields struct {
	BuffedBy   BuffDamage `json:"buffed"`
	DebuffedBy BuffDamage `json:"debuffed"`
	Buffed     int64      `json:"buffedDamage"`
	Debuffed   int64      `json:"debuffedDamage"`

	FADamage   int64   `json:"faDamage"`
	BADamage   int64   `json:"baDamage"`
	DeathTime  int64   `json:"deathTime"`
	DPSAverage []int64 `json:"dpsAverage"`
	DPSRolling []int64 `json:"dpsRolling"`

	Hits  int32 `json:"hits"`
	Casts int32 `json:"casts"`
	Crits int32 `json:"crits"`
	BA    int32 `json:"backAttacks"`
	FA    int32 `json:"frontAttacks"`
}

type StoredSkillFields struct {
	Icon         string     `json:"icon"`
	Hits         int32      `json:"hits"`
	Casts        int32      `json:"casts"`
	CastLog      []int32    `json:"castLog"`
	Crits        int32      `json:"crits"`
	MaxDamage    int64      `json:"maxDamage"`
	FA           int32      `json:"frontAttacks"`
	BA           int32      `json:"backAttacks"`
	FADamage     int64      `json:"faDamage"`
	BADamage     int64      `json:"baDamage"`
	TripodLevels TripodRows `json:"tripodLevels"`
	Buffed       BuffDamage `json:"buffed"`
	Debuffed     BuffDamage `json:"debuffed"`
}
