package process

import (
	"github.com/0fau/logs/pkg/process/meter"
	"github.com/cockroachdb/errors"
	"unicode"
)

func (p *Processor) Lint(enc *meter.Encounter) error {
	if !enc.DamageStats.Misc.Cleared {
		return errors.New("Uncleared")
	}

	if _, ok := enc.Entities[enc.LocalPlayer]; !ok {
		return errors.New("Missing data")
	}

	if enc.LocalPlayer == "You" && enc.Entities[enc.LocalPlayer].GearScore == 0 {
		return errors.New("Missing data")
	}

	players := make(map[string]struct{})
	for _, party := range enc.DamageStats.Misc.PartyInfo {
		for _, name := range party {
			if !ValidPlayerName(name) {
				return errors.New("Missing data")
			}

			if _, ok := players[name]; ok {
				return errors.New("Missing data")
			}

			players[name] = struct{}{}
		}
	}

	for name, entity := range enc.Entities {
		if entity.EntityType == "UNKNOWN" {
			return errors.New("Missing data")
		}

		if entity.EntityType != "PLAYER" {
			continue
		}

		if !ValidPlayerName(name) {
			return errors.New("Missing data")
		}

		if _, ok := players[name]; !ok {
			return errors.New("Missing data")
		}
	}

	return nil
}

func ValidPlayerName(name string) bool {
	if len(name) == 0 {
		return false
	}

	for i, r := range name {
		if i == 0 && unicode.IsLower(r) {
			return false
		}

		if !unicode.IsLetter(r) {
			return false
		}
	}

	return true
}
