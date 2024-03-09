package process

import "github.com/0fau/logs/pkg/process/meter"

func (p *Processor) Preprocess(raw *meter.Encounter) {
	if len(raw.DamageStats.Misc.PartyInfo) == 0 {
		var players []string
		for name, entity := range raw.Entities {
			if entity.EntityType != "PLAYER" {
				continue
			}

			players = append(players, name)
		}

		raw.DamageStats.Misc.PartyInfo = map[string][]string{"0": players}
	}

	if raw.Difficulty == "" {
		p.PopulateDifficulty(raw)
	}
}

func (p *Processor) PopulateDifficulty(raw *meter.Encounter) {
	switch raw.CurrentBossName {
	case "Sonavel", "Hanumatan":
		raw.Difficulty = "Normal"
	case "Saydon":
		if raw.DamageStats.TotalDamageDealt < 3000000000 {
			raw.Difficulty = "Normal"
		} else if raw.DamageStats.TotalDamageDealt > 4000000000 {
			raw.Difficulty = "Inferno"
		}
	case "Kakul":
		if raw.DamageStats.TotalDamageDealt < 3500000000 {
			raw.Difficulty = "Normal"
		} else {
			raw.Difficulty = "Inferno"
		}
	}
}
