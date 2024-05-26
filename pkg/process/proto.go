package process

//
//import (
//	"github.com/cockroachdb/errors"
//	"google.golang.org/protobuf/proto"
//
//	"github.com/0fau/logs/pkg/process/encproto"
//)
//
//func (enc *Encounter) SerializeData() ([]byte, error) {
//	buf := &encproto.EncounterData{
//		BuffCatalog:  make(map[string]*encproto.BuffInfo, len(enc.Data.BuffCatalog)),
//		SkillCatalog: make(map[string]*encproto.SkillInfo, len(enc.Data.SkillCatalog)),
//		Synergies:    make([]*encproto.BuffGroupsInfo, len(enc.Data.Synergies)),
//		Players:      make(map[string]*encproto.PlayerData, len(enc.Data.Players)),
//	}
//	for name, info := range enc.Data.BuffCatalog {
//		buf.BuffCatalog[name] = &encproto.BuffInfo{
//			Name:        info.Name,
//			Icon:        info.Icon,
//			Description: info.Description,
//			Category:    info.Category,
//			Set:         info.Set,
//		}
//		if info.Skill != nil {
//			buf.BuffCatalog[name].Skill = &encproto.BuffSkill{
//				Name:  info.Skill.Name,
//				Icon:  info.Skill.Icon,
//				Class: int32(info.Skill.Class),
//			}
//		}
//	}
//	for name, info := range enc.Data.SkillCatalog {
//		buf.SkillCatalog[name] = &encproto.SkillInfo{
//			Name: info.Name,
//			Icon: info.Icon,
//		}
//	}
//
//	for name, player := range enc.Data.Players {
//		buf.Players[name] = &encproto.PlayerData{
//			Damage: &encproto.PlayerDamage{
//				Crit:        player.Damage.Crit,
//				CritDamage:  player.Damage.CritDamage,
//				FrontAttack: player.Damage.FA,
//				BackAttack:  player.Damage.BA,
//			},
//		}
//	}
//
//	out, err := proto.Marshal(buf)
//	if err != nil {
//		return nil, errors.Wrap(err, "failed to encode encounter data")
//	}
//	return out, nil
//}
