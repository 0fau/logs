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
	Class   string `json:"class"`
	EntType string `json:"enttype"`
	Name    string `json:"name"`
	Damage  int64  `json:"damage"`
	Dps     int32  `json:"dps"`
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

		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	ret := make([]ReturnedEntity, len(entities))
	for i, ent := range entities {
		ret[i] = ReturnedEntity{
			Class:   ent.Class,
			EntType: ent.Enttype,
			Name:    ent.Name,
			Damage:  ent.Damage,
			Dps:     ent.Dps,
		}
	}
	c.JSON(http.StatusOK, ret)
}
