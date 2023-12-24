package process

import (
	"github.com/0fau/logs/pkg/process/meter"
	"github.com/cockroachdb/errors"
	"log"
	"unicode"
)

func (p *Processor) Lint(enc *meter.Encounter) error {
	if !enc.DamageStats.Misc.Cleared {
		return errors.New("Uncleared")
	}

	if _, ok := enc.Entities[enc.LocalPlayer]; !ok {
		log.Println("1")
		return errors.New("Missing data")
	}

	players := 0
	for name, player := range enc.Entities {
		if player.EntityType != "PLAYER" {
			continue
		}
		players++

		if !ValidPlayerName(name) {
			log.Println("2")
			return errors.New("Missing data")
		}
	}

	for _, party := range enc.DamageStats.Misc.PartyInfo {
		for _, name := range party {
			if !ValidPlayerName(name) {
				log.Println("3")
				return errors.New("Missing data")
			}
		}

		players -= len(party)
	}

	if len(enc.DamageStats.Misc.PartyInfo) > 0 && players != 0 {
		log.Println("4")
		return errors.New("Missing data")
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
