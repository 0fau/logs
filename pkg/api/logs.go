package api

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/0fau/logs/pkg/process/meter"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

type ReturnedEncounterShort struct {
	ID       int32  `json:"id"`
	Raid     string `json:"raid"`
	Date     int64  `json:"date"`
	Duration int32  `json:"duration"`
	Damage   int64  `json:"damage"`
}

type ReturnedEncounterDetails struct {
	Buffs     meter.BuffInfo   `json:"buffs"`
	Debuffs   meter.BuffInfo   `json:"debuffs"`
	HPLog     meter.HPLog      `json:"hpLog"`
	PartyInfo meter.PartyInfo  `json:"partyInfo"`
	Entities  []ReturnedEntity `json:"entities"`
}

type ReturnedEntity struct {
	Class      string           `json:"class"`
	EntType    string           `json:"enttype"`
	Name       string           `json:"name"`
	Damage     int64            `json:"damage"`
	FADamage   int64            `json:"faDamage"`
	BADamage   int64            `json:"baDamage"`
	Dps        int64            `json:"dps"`
	Dead       bool             `json:"dead"`
	DeathTime  int64            `json:"deathTime"`
	Skills     []ReturnedSkill  `json:"skills"`
	Buffed     meter.BuffDamage `json:"buffed"`
	Debuffed   meter.BuffDamage `json:"debuffed"`
	DPSAverage []int64          `json:"dpsAverage"`
	DPSRolling []int64          `json:"dpsRolling"`
}

type ReturnedSkill struct {
	SkillID     int32            `json:"id"`
	Casts       int32            `json:"casts"`
	CastLog     []int32          `json:"castLog"`
	Crits       int32            `json:"crits"`
	DPS         int64            `json:"dps"`
	Hits        int32            `json:"hits"`
	FADamage    int64            `json:"faDamage"`
	BADamage    int64            `json:"baDamage"`
	MaxDamage   int64            `json:"maxDamage"`
	TotalDamage int64            `json:"totalDamage"`
	TripodIndex meter.TripodRows `json:"tripodIndex"`
	TripodLevel meter.TripodRows `json:"tripodLevel"`
	Name        string           `json:"name"`
	Icon        string           `json:"icon"`
	Buffed      meter.BuffDamage `json:"buffed"`
	Debuffed    meter.BuffDamage `json:"debuffed"`
}

func (s *Server) recentLogs(c *gin.Context) {
	ctx := context.Background()

	var date time.Time
	if c.Query("past") != "" {
		num, err := strconv.ParseInt(c.Query("past"), 0, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		date = time.UnixMilli(num)
	}

	encs, err := s.conn.RecentEncounters(ctx, c.Query("user"), date)
	if err != nil {
		if strings.Contains(err.Error(), "scanning user uuid") {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		log.Println(errors.Wrap(err, "listing recent encounters"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ret := make([]ReturnedEncounterShort, len(encs))
	for i, enc := range encs {
		ret[i] = ReturnedEncounterShort{
			ID:       enc.ID,
			Raid:     enc.Raid,
			Date:     enc.Date.Time.UnixMilli(),
			Duration: enc.Duration,
			Damage:   enc.Damage,
		}
	}
	c.JSON(http.StatusOK, ret)
}

func (s *Server) logHandler(c *gin.Context) {
	if c.Param("log") == "" {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(c.Param("log"))
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	enc, err := s.conn.Queries.GetEncounter(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			log.Println(errors.Wrap(err, "fetching [encounter]"))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, ReturnedEncounterShort{
		ID:       enc.ID,
		Raid:     enc.Raid,
		Date:     enc.Date.Time.UnixMilli(),
		Duration: enc.Duration,
		Damage:   enc.Damage,
	})
}

func (s *Server) detailsHandler(c *gin.Context) {
	param := c.Param("log")
	id, err := strconv.Atoi(param)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var (
		entities []*sql.Entity
		fields   meter.StoredEncounterFields
		skills   []*sql.Skill
	)

	var wg sync.WaitGroup
	wg.Add(3)

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		defer wg.Done()
		var err error
		entities, err = s.conn.ListEntities(ctx, int32(id))
		if err != nil && !errors.Is(err, context.Canceled) {
			if errors.Is(err, pgx.ErrNoRows) {
				c.AbortWithStatus(http.StatusNotFound)
				return
			}
			cancel()

			log.Println(errors.Wrap(err, "listing entities"))
			return
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		fields, err = s.conn.Queries.GetFields(ctx, int32(id))
		if err != nil && !errors.IsAny(err, pgx.ErrNoRows, context.Canceled) {
			cancel()
			log.Println(errors.Wrap(err, "getting [encounter] buff info"))
			return
		}
	}()

	go func() {
		defer wg.Done()
		var err error
		skills, err = s.conn.ListSkills(ctx, int32(id))
		if err != nil && !errors.IsAny(err, pgx.ErrNoRows, context.Canceled) {
			cancel()
			log.Println(errors.Wrap(err, "listing [encounter] skills"))
			return
		}
	}()

	wg.Wait()

	if errors.Is(ctx.Err(), context.Canceled) {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	m := map[string]*ReturnedEntity{}
	ret := make([]ReturnedEntity, len(entities))
	for i, ent := range entities {
		ret[i] = ReturnedEntity{
			Class:      ent.Class,
			EntType:    ent.Enttype,
			Name:       ent.Name,
			Damage:     ent.Damage,
			FADamage:   ent.Fields.FADamage,
			BADamage:   ent.Fields.BADamage,
			Dps:        ent.Dps,
			DPSAverage: ent.Fields.DPSAverage,
			DPSRolling: ent.Fields.DPSRolling,
			Dead:       ent.Dead,
			DeathTime:  ent.Fields.DeathTime,
			Buffed:     ent.Fields.BuffedBy,
			Debuffed:   ent.Fields.DebuffedBy,
		}
		m[ent.Name] = &ret[i]
	}

	for _, skill := range skills {
		m[skill.Player].Skills = append(
			m[skill.Player].Skills,
			ReturnedSkill{
				SkillID:     skill.SkillID,
				Casts:       skill.Fields.Casts,
				CastLog:     skill.Fields.CastLog,
				Crits:       skill.Fields.Crits,
				Buffed:      skill.Fields.Buffed,
				Debuffed:    skill.Fields.Debuffed,
				Hits:        skill.Fields.Hits,
				FADamage:    skill.Fields.FADamage,
				BADamage:    skill.Fields.BADamage,
				MaxDamage:   skill.Fields.MaxDamage,
				TripodLevel: skill.Fields.TripodLevels,
				TripodIndex: skill.Tripods,
				Icon:        skill.Fields.Icon,
				DPS:         skill.Dps,
				TotalDamage: skill.Damage,
				Name:        skill.Name,
			},
		)
	}

	c.JSON(http.StatusOK, ReturnedEncounterDetails{
		Buffs:     fields.Buffs,
		Debuffs:   fields.Debuffs,
		HPLog:     fields.HPLog,
		PartyInfo: fields.PartyInfo,
		Entities:  ret,
	})
}
