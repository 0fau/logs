package api

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ReturnedEncounterShort struct {
	User        ReturnedUser `json:"user"`
	ID          int32        `json:"id"`
	Difficulty  string       `json:"difficulty"`
	Boss        string       `json:"boss"`
	Date        int64        `json:"date"`
	Duration    int32        `json:"duration"`
	LocalPlayer string       `json:"localPlayer"`
	structs.EncounterHeader
}

type ReturnedEncounter struct {
	ReturnedEncounterShort
	Data structs.EncounterData `json:"data"`
}

type ReturnedEncounterShorts struct {
	Encounters []ReturnedEncounterShort `json:"encounters"`
	More       bool                     `json:"more"`
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
		date = time.UnixMilli(num).UTC()
		log.Println(date)
	}

	var id int32
	if c.Query("id") != "" {
		num, err := strconv.Atoi(c.Query("id"))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		id = int32(num)
	}

	user := ""
	if c.Query("scope") == "roster" || c.Query("scope") == "friends" {
		sesh := sessions.Default(c)
		val := sesh.Get("user")
		if val == nil {
			c.JSON(http.StatusUnauthorized, []struct{}{})
			return
		}
		u := val.(*SessionUser)
		user = u.ID
	}

	encs, err := s.conn.RecentEncounters(ctx, user, id, date, c.Query("scope") == "friends")
	if err != nil {
		if strings.Contains(err.Error(), "scanning user uuid") {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}

		log.Println(errors.Wrap(err, "listing recent encounters"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	n, more := len(encs), false
	if n > 5 {
		n = 5
		more = true
	}

	ret := make([]ReturnedEncounterShort, n)
	for i, enc := range encs[:n] {
		ret[i] = ReturnedEncounterShort{
			ID:              enc.ID,
			Difficulty:      enc.Difficulty,
			Boss:            enc.Boss,
			Date:            enc.Date.Time.UnixMilli(),
			Duration:        enc.Duration,
			LocalPlayer:     enc.LocalPlayer,
			EncounterHeader: enc.Header,
		}
	}
	c.JSON(http.StatusOK, ReturnedEncounterShorts{
		Encounters: ret,
		More:       more,
	})
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

	uploadedBy, _ := enc.UploadedBy.Value()
	username := ""
	if u, _ := enc.Username.Value(); u != nil {
		username = u.(string)
	}
	avatar, _ := enc.Avatar.Value()

	c.JSON(http.StatusOK, ReturnedEncounter{
		ReturnedEncounterShort: ReturnedEncounterShort{
			User: ReturnedUser{
				ID:         uploadedBy.(string),
				DiscordTag: enc.DiscordTag,
				DiscordID:  enc.DiscordID,
				Username:   username,
				Avatar:     avatar.(string),
			},
			ID:              enc.ID,
			Difficulty:      enc.Difficulty,
			Boss:            enc.Boss,
			Date:            enc.Date.Time.UnixMilli(),
			Duration:        enc.Duration,
			LocalPlayer:     enc.LocalPlayer,
			EncounterHeader: enc.Header,
		},
		Data: enc.Data,
	})
}
