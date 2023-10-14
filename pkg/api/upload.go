package api

import (
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/0fau/logs/pkg/process/meter"
	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	json "github.com/goccy/go-json"
	"github.com/jackc/pgx/v5"
	"github.com/thanhpk/randstr"
	"io"
	"log"
	"net/http"
)

func (s *Server) generateToken(c *gin.Context) {
	sesh := sessions.Default(c)

	var u *SessionUser
	if val := sesh.Get("user"); val != nil {
		u = val.(*SessionUser)
	} else {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	token := randstr.String(64)
	if err := s.conn.SetUserAccessToken(ctx, u.ID, token); err != nil {
		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, token)
}

func hasRole(user *sql.User, role string) bool {
	for _, has := range user.Roles {
		if has == role {
			return true
		}
	}
	return false
}

func (s *Server) uploadHandler(c *gin.Context) {
	token := c.GetHeader("access_token")
	if token == "" {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	user, err := s.conn.UserByAccessToken(ctx, token)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			c.AbortWithStatus(http.StatusUnauthorized)
		} else {
			log.Println(errors.WithStack(err))
			c.AbortWithStatus(http.StatusInternalServerError)
		}
		return
	}

	if !hasRole(user, "trusted") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer c.Request.Body.Close()

	enc := &meter.Encounter{}
	if err := json.Unmarshal(raw, &enc); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := s.processor.Lint(enc); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if _, err := s.processor.Save(ctx, s.conn, user.ID, enc); err != nil {
		log.Println(errors.Wrap(err, "saving encounter"))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
}
