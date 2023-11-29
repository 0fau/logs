package api

import (
	"cmp"
	"context"
	"fmt"
	"github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"log"
	"net/http"
	"slices"
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
	Anonymized  bool         `json:"anonymized"`
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

func (enc *ReturnedEncounter) Anonymize(order map[string]string) {
	for name, player := range enc.Data.Players {
		enc.Data.Players[order[name]] = player
		delete(enc.Data.Players, name)
	}
	enc.ReturnedEncounterShort.Anonymize(order)
}

func (enc *ReturnedEncounterShort) Order() map[string]string {
	players := make([]string, 0, len(enc.Players))
	for name := range enc.Players {
		players = append(players, name)
	}
	slices.SortFunc(players, func(a, b string) int {
		return cmp.Compare(enc.Players[b].Damage, enc.Players[a].Damage)
	})
	m := map[string]string{}
	for i, player := range players {
		m[player] = fmt.Sprintf("#%d", i+1)
	}
	return m
}

func (enc *ReturnedEncounterShort) Anonymize(order map[string]string) {
	parties := make([][]string, 0, len(enc.Parties))
	for _, party := range enc.Parties {
		anon := make([]string, 0, len(party))
		for _, player := range party {
			anon = append(anon, order[player])
		}
		parties = append(parties, anon)
	}
	enc.Parties = parties

	for name, player := range enc.Players {
		enc.Players[order[name]] = player
		delete(enc.Players, name)
	}

	enc.LocalPlayer = order[enc.LocalPlayer]
	enc.User = ReturnedUser{}
	enc.Anonymized = true
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
		username := ""
		if name, _ := enc.Username.Value(); name != nil {
			username = name.(string)
		}
		short := ReturnedEncounterShort{
			User: ReturnedUser{
				DiscordTag: enc.DiscordTag,
				Username:   username,
			},
			ID:              enc.ID,
			Difficulty:      enc.Difficulty,
			Boss:            enc.Boss,
			Date:            enc.Date.Time.UnixMilli(),
			Duration:        enc.Duration,
			LocalPlayer:     enc.LocalPlayer,
			EncounterHeader: enc.Header,
		}
		//order := short.Order()
		//short.Anonymize(order)

		ret[i] = short
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

	full := ReturnedEncounter{
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
	}
	//order := full.Order()
	//full.Anonymize(order)

	c.JSON(http.StatusOK, full)
}
