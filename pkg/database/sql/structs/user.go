package structs

type UserSettings struct {
	SkipLanding       bool   `json:"skip_landing"`
	ProfileVisibility string `json:"profile_visibility"`
	LogVisibility     string `json:"log_visibility"`
}
