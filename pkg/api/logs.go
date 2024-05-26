package api

import (
	"cmp"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"slices"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/goccy/go-json"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/0fau/logs/pkg/database/sql"
	structs2 "github.com/0fau/logs/pkg/process/structs"
	"github.com/0fau/logs/pkg/query"
)

type ReturnedEncounterShort struct {
	Uploader    *ReturnedUser                 `json:"uploader"`
	ID          int32                         `json:"id"`
	Difficulty  string                        `json:"difficulty"`
	Boss        string                        `json:"boss"`
	Date        int64                         `json:"date"`
	Duration    int32                         `json:"duration"`
	LocalPlayer string                        `json:"localPlayer"`
	Anonymized  bool                          `json:"anonymized"`
	Actual      bool                          `json:"actual"`
	Place       int32                         `json:"place"`
	Visibility  *structs2.EncounterVisibility `json:"visibility"`
	structs2.EncounterHeader
}

type ReturnedEncounter struct {
	ReturnedEncounterShort
	Thumbnail bool                   `json:"thumbnail"`
	Data      structs2.EncounterData `json:"data"`
}

type ReturnedEncounterShorts struct {
	Encounters []ReturnedEncounterShort `json:"encounters"`
	More       bool                     `json:"more"`
}

func (enc *ReturnedEncounter) Anonymize(order map[string]string, character string, visibility int) {
	m := make(map[string]structs2.PlayerData, len(enc.Data.Players))
	for name, player := range enc.Data.Players {
		if visibility == structs2.ShowSelf && name == character {
			m[name] = player
		} else {
			m[order[name]] = player
		}
	}
	enc.Data.Players = m
	enc.ReturnedEncounterShort.Anonymize(order, character, visibility)
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
		m[player] = fmt.Sprintf("%s #%d", enc.Players[player].Class, i+1)
	}
	return m
}

func (enc *ReturnedEncounterShort) Anonymize(order map[string]string, character string, visibility int) {
	parties := make([][]string, 0, len(enc.Parties))
	for _, party := range enc.Parties {
		anon := make([]string, 0, len(party))
		for _, player := range party {
			if visibility == structs2.ShowSelf && player == character {
				anon = append(anon, player)
			} else {
				anon = append(anon, order[player])
			}
		}
		parties = append(parties, anon)
	}
	enc.Parties = parties

	m := make(map[string]structs2.PlayerHeader, len(enc.Players))
	for name, player := range enc.Players {
		if visibility == structs2.ShowSelf && name == character {
			m[name] = player
		} else {
			player.Name = order[name]
			m[order[name]] = player
		}
	}
	enc.Players = m
	if (visibility == structs2.ShowSelf && enc.LocalPlayer != character) || visibility == structs2.HideNames || visibility == structs2.UnsetNames {
		enc.LocalPlayer = order[enc.LocalPlayer]
		enc.Anonymized = true
	}

	if visibility != structs2.ShowSelf {
		enc.Uploader = nil
	}
}

func (s *Server) logs(c *gin.Context) {
	params := &query.Params{
		Order:     c.Query("order"),
		Scope:     c.Query("scope"),
		GearScore: c.Query("gear_score"),
	}

	if !slices.Contains([]string{"arkesia", "friends", "roster"}, c.Query("scope")) {
		params.Scope = "arkesia"
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

	if c.Query("past_place") != "" {
		num, err := strconv.Atoi(c.Query("past_place"))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		cast := int32(num)
		params.PastPlace = &cast
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

	trusted := hasRoles(roles, "trusted")
	admin := hasRoles(roles, "admin")

	params.User = uuid
	params.Privileged = trusted || admin
	params.Admin = admin

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
			Actual:          enc.LocalPlayer == enc.FocusedPlayer || enc.FocusedPlayer == "",
			EncounterHeader: enc.Header,
			Place:           enc.Place,
		}
		if enc.FocusedPlayer != "" {
			short.LocalPlayer = enc.FocusedPlayer
		}

		uploaderID, _ := enc.UploadedBy.Value()

		visibility := structs2.UnsetNames
		if enc.Visibility != nil && enc.Visibility.Names != structs2.UnsetNames {
			visibility = enc.Visibility.Names
		} else if enc.Uploader.LogVisibility != nil {
			visibility = enc.Uploader.LogVisibility.Names
		}

		showNames := visibility == structs2.ShowNames
		if !(showNames || (u != nil &&
			(u.ID == uploaderID || params.Scope == "friends" ||
				(!(enc.Visibility != nil && (enc.Visibility.Names == structs2.HideNames || enc.Visibility.Names == structs2.ShowSelf)) &&
					hasRoles(roles, "trusted")) || hasRoles(roles, "admin")))) {
			order := short.Order()
			short.Anonymize(order, enc.LocalPlayer, visibility)
		}
		short.Visibility = &structs2.EncounterVisibility{Names: visibility}

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
	if enc.Username != nil {
		username = *enc.Username
	}

	if enc.Private && (u == nil || u.ID != uploadedBy) && !hasRoles(roles, "admin") {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	friends, err := s.conn.Queries.AreFriends(ctx, sql.AreFriendsParams{
		User1: uuid, User2: enc.UploadedBy,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		log.Println(errors.Wrap(err, "fetching rows"))
		c.AbortWithStatus(http.StatusInternalServerError)
	}

	full := ReturnedEncounter{
		ReturnedEncounterShort: ReturnedEncounterShort{
			Uploader: &ReturnedUser{
				ID:         uploadedBy.(string),
				DiscordTag: enc.DiscordTag,
				Username:   username,
				Avatar:     enc.Avatar != "",
			},
			ID:              enc.ID,
			Difficulty:      enc.Difficulty,
			Boss:            enc.Boss,
			Date:            enc.Date.Time.UnixMilli(),
			Duration:        enc.Duration,
			LocalPlayer:     enc.LocalPlayer,
			Actual:          true,
			EncounterHeader: enc.Header,
		},
		Thumbnail: enc.Thumbnail,
	}

	if enc.Data != nil {
		full.Data = *enc.Data
	} else {
		b, err := s.s3.FetchEncounter(ctx, enc.ID)
		if err != nil {
			log.Println(errors.Wrap(err, "fetching encounter data from s3"))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if err := json.Unmarshal(b, &full.Data); err != nil {
			log.Println(errors.Wrap(err, "unmarshaling encounter data"))
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	visibility := structs2.UnsetNames
	if enc.Visibility != nil && enc.Visibility.Names != structs2.UnsetNames {
		visibility = enc.Visibility.Names
	} else if enc.LogVisibility != nil {
		visibility = enc.LogVisibility.Names
	}

	showNames := visibility == structs2.ShowNames
	if !(showNames || (u != nil &&
		(u.ID == uploadedBy || friends || (!(enc.Visibility != nil && (enc.Visibility.Names == structs2.HideNames ||
			enc.Visibility.Names == structs2.ShowSelf)) && hasRoles(roles, "trusted")) || hasRoles(roles, "admin")))) {
		order := full.Order()
		full.Anonymize(order, enc.LocalPlayer, visibility)
	}
	full.ReturnedEncounterShort.Visibility = &structs2.EncounterVisibility{Names: visibility}

	c.JSON(http.StatusOK, full)
}

func (s *Server) shortLogHandler(c *gin.Context) {
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
	enc, err := s.conn.Queries.GetEncounterShort(ctx, int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			log.Println(errors.Wrap(err, "fetching encounter"))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	short := ReturnedEncounterShort{
		ID:              enc.ID,
		Difficulty:      enc.Difficulty,
		Boss:            enc.Boss,
		Date:            enc.Date.Time.UnixMilli(),
		Duration:        enc.Duration,
		LocalPlayer:     enc.LocalPlayer,
		EncounterHeader: enc.Header,
	}

	order := short.Order()
	short.Anonymize(order, "", structs2.HideNames)

	c.JSON(http.StatusOK, short)
}

type UpdateLogSettingsParams struct {
	Log int32 `uri:"log" binding:"required"`
}

type UpdateLogSettingsBody struct {
	Names int `json:"names"`
}

func (s *Server) updateLogSettings(c *gin.Context) {
	sesh := sessions.Default(c)

	var u *SessionUser
	if val := sesh.Get("user"); val != nil {
		u = val.(*SessionUser)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	params := UpdateLogSettingsParams{}
	if err := c.ShouldBindUri(&params); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx := context.Background()
	row, err := s.conn.Queries.GetEncounterVisibility(ctx, params.Log)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			log.Println(errors.Wrap(err, "fetching encounter visibility"))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	var uuid string
	if uploader, err := row.UploadedBy.Value(); err == nil && uploader != nil {
		uuid = uploader.(string)
	}

	if uuid != u.ID {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if row.Visibility == nil {
		row.Visibility = &structs2.EncounterVisibility{}
	}

	body := UpdateLogSettingsBody{}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	if body.Names != row.Visibility.Names {
		row.Visibility.Names = body.Names
	} else {
		return
	}

	if err := s.conn.Queries.UpdateEncounterVisibility(ctx, sql.UpdateEncounterVisibilityParams{
		ID: params.Log, Visibility: row.Visibility,
	}); err != nil {
		log.Println(errors.Wrap(err, "updating encounter visibility"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
