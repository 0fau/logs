package api

import (
	"context"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"strconv"
)

type ReturnedEncounterShort struct {
	ID       int    `json:"id"`
	Raid     string `json:"raid"`
	Date     int64  `json:"date"`
	Duration int    `json:"duration"`
	Damage   int64  `json:"damage"`
}

type ReturnedEntity struct {
	Class   string          `json:"class"`
	EntType string          `json:"enttype"`
	Name    string          `json:"name"`
	Damage  int64           `json:"damage"`
	Dps     int64           `json:"dps"`
	Skills  []ReturnedSkill `json:"skills"`
}

type ReturnedSkill struct {
	SkillID     int32  `json:"id"`
	Casts       int32  `json:"casts"`
	Crits       int32  `json:"crits"`
	DPS         int64  `json:"dps"`
	Hits        int32  `json:"hits"`
	MaxDamage   int64  `json:"maxDamage"`
	TotalDamage int64  `json:"totalDamage"`
	Name        string `json:"name"`
}

func (s *Server) recentLogs(c *gin.Context) {
	ctx := context.Background()
	encs, err := s.conn.RecentEncounters(ctx)
	if err != nil {
		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ret := make([]ReturnedEncounterShort, len(encs))
	for i, enc := range encs {
		ret[i] = ReturnedEncounterShort{
			ID:       int(enc.ID),
			Raid:     enc.Raid,
			Date:     enc.Date.Time.UnixMilli(),
			Duration: int(enc.Duration),
			Damage:   enc.TotalDamageDealt,
		}
	}
	c.JSON(http.StatusOK, ret)
}

func (s *Server) logHandler(c *gin.Context) {
	ctx := context.Background()
	param := c.Param("log")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	entities, err := s.conn.ListEntities(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}

		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	skills, err := s.conn.ListSkills(ctx, int32(id))
	if err != nil && errors.Is(err, pgx.ErrNoRows) {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	m := map[string]*ReturnedEntity{}
	ret := make([]ReturnedEntity, len(entities))
	for i, ent := range entities {
		ret[i] = ReturnedEntity{
			Class:   ent.Class,
			EntType: ent.Enttype,
			Name:    ent.Name,
			Damage:  ent.Damage,
			Dps:     ent.Dps,
		}
		m[ent.Name] = &ret[i]
	}

	for _, skill := range skills {
		m[skill.Player].Skills = append(
			m[skill.Player].Skills,
			ReturnedSkill{
				SkillID:     skill.SkillID,
				Casts:       skill.Casts,
				Crits:       skill.Crits,
				DPS:         skill.Dps,
				Hits:        skill.Hits,
				MaxDamage:   skill.MaxDamage,
				TotalDamage: skill.TotalDamage,
				Name:        skill.Name,
			},
		)
	}

	c.JSON(http.StatusOK, ret)
}
