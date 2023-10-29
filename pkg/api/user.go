package api

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"log"
	"net/http"
	"unicode"
)

func (s *Server) userHandler(c *gin.Context) {
	ctx := context.Background()
	user, err := s.conn.Queries.GetUser(ctx, c.Param("user"))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusNotFound)
		} else {
			log.Println(errors.Wrap(err, "get user by name"))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	uuid, _ := user.ID.Value()
	c.JSON(http.StatusOK, ReturnedUser{
		ID:       uuid.(string),
		Username: c.Param("user"),
	})
}

type ReturnedSettings struct {
	HasToken bool `json:"hasToken"`
	structs.UserSettings
}

func (s *Server) settingsHandler(c *gin.Context) {
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
	user, err := s.conn.Queries.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			log.Println(errors.Wrap(err, "get user by name"))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	log.Println(user.AccessToken.Valid)

	c.JSON(http.StatusOK, &ReturnedSettings{
		HasToken:     user.AccessToken.Valid,
		UserSettings: user.Settings,
	})
}

func (s *Server) setUsername(c *gin.Context) {
	sesh := sessions.Default(c)
	var u *SessionUser
	if val := sesh.Get("user"); val != nil {
		u = val.(*SessionUser)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if u.Username == c.Query("username") {
		return
	}

	if c.Query("username") == "" || len(c.Query("username")) > 16 {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	for _, r := range c.Query("username") {
		if !(unicode.IsLetter(r) || unicode.IsDigit(r)) {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
	}

	id := pgtype.UUID{}
	if err := id.Scan(u.ID); err != nil {
		log.Println(errors.Wrap(err, "scanning pgtype.UUID id"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	username := pgtype.Text{}
	if err := username.Scan(c.Query("username")); err != nil {
		log.Println(errors.Wrap(err, "scanning pgtype.Text username"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	if err := s.conn.Queries.SetUsername(ctx, sql.SetUsernameParams{
		ID:       id,
		Username: username,
	}); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	u.Username = c.Query("username")
	sesh.Set("user", u)
	if err := sesh.Save(); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
