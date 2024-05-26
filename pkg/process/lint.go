package process

import (
	"cmp"
	"slices"
	"strings"
	"unicode"

	"github.com/cockroachdb/errors"

	"github.com/0fau/logs/pkg/process/meter"
)

func (p *Processor) Lint(enc *meter.Encounter) (error, int) {
	if !enc.DamageStats.Misc.Cleared {
		return errors.New("Uncleared"), 1
	}

	if _, ok := enc.Entities[enc.LocalPlayer]; !ok {
		return errors.New("Missing data"), 2
	}

	if enc.LocalPlayer == "You" && enc.Entities[enc.LocalPlayer].GearScore == 0 {
		return errors.New("Missing data"), 3
	}

	players := make(map[string]struct{})
	for _, party := range enc.DamageStats.Misc.PartyInfo {
		for _, name := range party {
			if !ValidPlayerName(name) {
				return errors.New("Missing data"), 4
			}

			if _, ok := players[name]; ok {
				return errors.New("Missing data"), 5
			}

			players[name] = struct{}{}
		}
	}

	hps := make([][]meter.HPLogEntry, 0, len(enc.DamageStats.Misc.HPLog))
	for boss, hp := range enc.DamageStats.Misc.HPLog {
		if len(hp) == 0 {
			continue
		}

		if _, ok := RaidLookup[boss]; ok {
			hps = append(hps, hp)
			continue
		}

		if _, ok := GuardianLookup[boss]; ok {
			hps = append(hps, hp)
			continue
		}
	}
	if len(hps) != 0 {
		hp := slices.MinFunc(hps, func(a, b []meter.HPLogEntry) int {
			return cmp.Compare(a[0].Time, b[0].Time)
		})
		var foundHP bool
		for i := 0; i < 10 && i < len(hp); i++ {
			if hp[i].P > 0.98 {
				foundHP = true
				break
			}
		}
		if !foundHP {
			return errors.New("Missing data"), 6
		}
	}

	for name, entity := range enc.Entities {
		if entity.EntityType == "UNKNOWN" || entity.EntityType == "NPC" {
			return errors.New("Missing data"), 7
		}

		if entity.EntityType != "PLAYER" {
			continue
		}

		if !ValidPlayerName(name) {
			return errors.New("Missing data"), 8
		}

		if !p.ValidPlayer(entity) {
			return errors.New("Missing data"), 10
		}

		if _, ok := players[name]; !ok {
			return errors.New("Missing data"), 9
		}
	}

	return nil, 0
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

func (p *Processor) ValidPlayer(ent meter.Entity) bool {
	for name, skill := range ent.Skills {
		lookup, ok := p.skills[name]
		if !ok {
			continue
		}

		if !strings.HasPrefix(skill.Icon, "battle_item") && lookup.ClassID != 0 && lookup.ClassID != ent.ClassId {
			return false
		}
	}

	return true
}
