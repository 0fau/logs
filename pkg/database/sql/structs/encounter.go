package structs

type PlayerHeader struct {
	Name      string    `json:"name"`
	Class     string    `json:"class"`
	Damage    int64     `json:"damage"`
	DPS       int64     `json:"dps"`
	Alive     bool      `json:"alive"`
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
	PartyDamage [][]int64
	PartyBuffs  [][]string
	PartyBuffed [][][]int64

	PlayerDamage [][]int64
	PlayerBuff   [][]int64

	PartySelfBuff   [][]string
	PartySelfBuffed [][][]int64

	PlayerSelfBuff   [][]string
	PlayerSelfBuffed [][][]int64
}

type EncounterSettings struct {
	Visibility []string
}
