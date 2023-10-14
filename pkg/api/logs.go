package api

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/0fau/logs/pkg/process/meter"
	"github.com/cockroachdb/errors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ReturnedEncounterShort struct {
	ID       int32  `json:"id"`
	Boss     string `json:"boss"`
	Date     int64  `json:"date"`
	Duration int32  `json:"duration"`
	structs.EncounterHeader
	structs.EncounterData
}

type ReturnedEncounterDetails struct {
	Buffs     meter.BuffInfo   `json:"buffs"`
	Debuffs   meter.BuffInfo   `json:"debuffs"`
	HPLog     meter.HPLog      `json:"hpLog"`
	PartyInfo meter.PartyInfo  `json:"partyInfo"`
	Entities  []ReturnedEntity `json:"entities"`

	structs.EncounterData
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
			ID:              enc.ID,
			Boss:            enc.Boss,
			Date:            enc.Date.Time.UnixMilli(),
			Duration:        enc.Duration,
			EncounterHeader: enc.Header,
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
			log.Println(errors.Wrap(err, "fetching encounter"))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	c.JSON(http.StatusOK, ReturnedEncounterShort{
		ID:              enc.ID,
		Boss:            enc.Boss,
		Date:            enc.Date.Time.UnixMilli(),
		Duration:        enc.Duration,
		EncounterHeader: enc.Header,
	})
}
