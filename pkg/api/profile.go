package api

import (
	"context"
	"log"
	"net/http"
	"slices"

	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgtype"
)

type ReturnedCharacter struct {
	Character string  `json:"character"`
	Class     string  `json:"class"`
	GearScore float64 `json:"gearScore"`
	LogCount  int64   `json:"logCount"`
}

type ReturnedProfile struct {
	Characters []ReturnedCharacter `json:"characters"`
}

func (s *Server) profileHandler(c *gin.Context) {
	sesh := sessions.Default(c)
	var u *SessionUser
	if val := sesh.Get("user"); val != nil {
		u = val.(*SessionUser)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	id := pgtype.UUID{}
	if err := id.Scan(u.ID); err != nil {
		log.Println(errors.Wrap(err, "scanning pgtype.UUID id"))
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	characters, err := s.conn.Queries.GetRosterStatsByID(ctx, id)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var ret []ReturnedCharacter
	for _, character := range characters {
		ret = append(ret, ReturnedCharacter{
			Character: character.Name,
			Class:     character.Class,
			GearScore: character.GearScore,
			LogCount:  character.Count,
		})
	}

	c.JSON(http.StatusOK, ret)
}

type UserProfileParams struct {
	Username string `uri:"username" binding:"required"`
}

func (s *Server) userProfileHandler(c *gin.Context) {
	sesh := sessions.Default(c)

	var roles []string
	if val := sesh.Get("roles"); val != nil {
		roles = val.([]string)
	}

	if !slices.Contains(roles, "admin") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var params UserProfileParams
	if err := c.ShouldBindUri(&params); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	ctx := context.Background()
	characters, err := s.conn.Queries.GetRosterStats(ctx, params.Username)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var ret []ReturnedCharacter
	for _, character := range characters {
		ret = append(ret, ReturnedCharacter{
			Character: character.Name,
			Class:     character.Class,
			GearScore: character.GearScore,
			LogCount:  character.Count,
		})
	}

	c.JSON(http.StatusOK, ret)
}
