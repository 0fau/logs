package api

import (
	"cmp"
	"context"
	"fmt"
	"github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/0fau/logs/pkg/query"
	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"
	"strings"
	"time"
)

type ReturnedEncounterShort struct {
	Uploader    ReturnedUser `json:"uploader"`
	ID          int32        `json:"id"`
	Difficulty  string       `json:"difficulty"`
	Boss        string       `json:"boss"`
	Date        int64        `json:"date"`
	Duration    int32        `json:"duration"`
	LocalPlayer string       `json:"localPlayer"`
	Place       int32        `json:"place"`
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
	m := make(map[string]structs.PlayerData, len(enc.Data.Players))
	for name, player := range enc.Data.Players {
		m[order[name]] = player
	}
	enc.Data.Players = m
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

	m := make(map[string]structs.PlayerHeader, len(enc.Players))
	for name, player := range enc.Players {
		m[order[name]] = player
	}
	enc.Players = m

	enc.LocalPlayer = order[enc.LocalPlayer]
	enc.Uploader = ReturnedUser{}
	enc.Anonymized = true
}

func (s *Server) logs(c *gin.Context) {
	params := &query.Params{
		Order: c.Query("order"),
		Scope: c.Query("scope"),
	}

	if c.Query("past_id") != "" {
		num, err := strconv.Atoi(c.Query("past_id"))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		cast := int32(num)
		params.PastID = &cast
	}

	if c.Query("past_field") != "" {
		num, err := strconv.ParseInt(c.Query("past_field"), 10, 64)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		params.PastField = &num
	}

	body, err := io.ReadAll(c.Request.Body)
	defer func() {
		if err := c.Request.Body.Close(); err != nil {
			log.Println(errors.Wrap(err, "couldn't close request body"))
		}
	}()
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if len(body) != 0 {
		if err := json.Unmarshal(body, &params.Selections); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	sesh := sessions.Default(c)
	val := sesh.Get("user")

	var u *SessionUser
	if val != nil {
		u = val.(*SessionUser)
	}

	var uuid pgtype.UUID
	if u != nil {
		if err := uuid.Scan(u.ID); err != nil {
			log.Println(errors.Wrap(err, "scanning session user uuid"))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	roles, err := s.conn.Queries.GetRoles(context.Background(), uuid)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Println(errors.Wrap(err, "fetching rows"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	params.User = uuid
	params.Privileged = hasRoles(roles, "admin", "trusted")

	encs, err := query.Query(s.conn, params)
	if err != nil {
		log.Println(errors.Wrap(err, "query"))
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
		short := ReturnedEncounterShort{
			ID:              enc.ID,
			Difficulty:      enc.Difficulty,
			Boss:            enc.Boss,
			Date:            enc.Date.Time.UnixMilli(),
			Duration:        enc.Duration,
			LocalPlayer:     enc.LocalPlayer,
			EncounterHeader: enc.Header,
			Place:           enc.Place,
		}

		uploaderID, _ := enc.UploadedBy.Value()
		if u == nil || uploaderID != u.ID && !hasRoles(roles, "admin", "trusted") {
			order := short.Order()
			short.Anonymize(order)
		}

		ret[i] = short
	}
	c.JSON(http.StatusOK, ReturnedEncounterShorts{
		Encounters: ret,
		More:       more,
	})
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

	sesh := sessions.Default(c)
	val := sesh.Get("user")

	var u *SessionUser
	if val != nil {
		u = val.(*SessionUser)
	}

	focus := ""
	if c.Query("scope") == "roster" || c.Query("scope") == "friends" {
		if u == nil {
			c.JSON(http.StatusUnauthorized, []struct{}{})
			return
		}

		focus = u.ID
	}

	var uuid pgtype.UUID
	if u != nil {
		if err := uuid.Scan(u.ID); err != nil {
			log.Println(errors.Wrap(err, "scanning session user uuid"))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	roles, err := s.conn.Queries.GetRoles(ctx, uuid)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Println(errors.Wrap(err, "fetching rows"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	encs, err := s.conn.RecentEncounters(ctx, focus, id, date, c.Query("scope") == "friends")
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
		short := ReturnedEncounterShort{
			ID:              enc.ID,
			Difficulty:      enc.Difficulty,
			Boss:            enc.Boss,
			Date:            enc.Date.Time.UnixMilli(),
			Duration:        enc.Duration,
			LocalPlayer:     enc.LocalPlayer,
			EncounterHeader: enc.Header,
		}

		uploaderID, _ := enc.UploadedBy.Value()
		if u == nil || uploaderID != u.ID && !hasRoles(roles, "admin", "trusted") {
			order := short.Order()
			short.Anonymize(order)
		}

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

	sesh := sessions.Default(c)
	val := sesh.Get("user")

	var u *SessionUser
	if val != nil {
		u = val.(*SessionUser)
	}

	var uuid pgtype.UUID
	if u != nil {
		if err := uuid.Scan(u.ID); err != nil {
			log.Println(errors.Wrap(err, "scanning session user uuid"))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	roles, err := s.conn.Queries.GetRoles(ctx, uuid)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Println(errors.Wrap(err, "fetching rows"))
		c.AbortWithStatus(http.StatusInternalServerError)
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
			Uploader: ReturnedUser{
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

	if u == nil || u.ID != uploadedBy && !hasRoles(roles, "admin", "trusted") {
		order := full.Order()
		full.Anonymize(order)
	}

	c.JSON(http.StatusOK, full)
}
