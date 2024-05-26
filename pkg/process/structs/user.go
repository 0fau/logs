package structs

type UserSettings struct {
	SkipLanding       bool   `json:"skip_landing"`
	ProfileVisibility string `json:"profile_visibility"`
	LogVisibility     string `json:"log_visibility"`
}

const (
	UnsetNames = iota
	ShowNames
	ShowSelf
	HideNames
)

type EncounterVisibility struct {
	Names int `json:"names"`
}
