package structs

type PlayerHeader struct {
	Name   string `json:"name"`
	Class  string `json:"class"`
	Damage int64  `json:"damage"`
	DPS    int64  `json:"dps"`
	Alive  bool   `json:"alive"`
}

type EncounterHeader struct {
	Players map[string]PlayerHeader `json:"players"`
	Parties [][]string              `json:"parties"`
	Damage  int64                   `json:"damage"`
	Cleared bool                    `json:"cleared"`
}

type EncounterData struct {
}

type EncounterSettings struct {
	Visibility []string
}
