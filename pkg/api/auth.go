package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"io"
	"log"
	"net/http"
	"time"
)

type DiscordUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

type SessionUser struct {
	ID          string
	DiscordName string
	CreatedAt   time.Time
	Roles       []string
}

type ReturnedUser struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
}

func redirectLoggedIn(c *gin.Context) {
	c.Redirect(http.StatusFound, "/logs")
}

func (s *Server) oauth2(c *gin.Context) {
	sesh := sessions.Default(c)
	if user := sesh.Get("user"); user != nil {
		redirectLoggedIn(c)
		return
	}

	state := randstr.String(32)
	url := s.conf.OAuth2.AuthCodeURL(state)
	sesh.Set("oauth_state", state)
	sesh.Save()

	c.Redirect(http.StatusFound, url)
}

func (s *Server) oauth2Redirect(c *gin.Context) {
	sesh := sessions.Default(c)
	if user := sesh.Get("user"); user != nil {
		redirectLoggedIn(c)
		return
	}

	state := sesh.Get("oauth_state")
	if state == nil {
		c.Status(http.StatusUnauthorized)
		return
	}
	sesh.Delete("oauth_state")
	if err := sesh.Save(); err != nil {
		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if state.(string) != c.Query("state") {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	ctx := context.Background()
	token, err := s.conf.OAuth2.Exchange(ctx, c.Query("code"))
	if err != nil {
		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client := s.conf.OAuth2.Client(ctx, token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	u := DiscordUser{}
	if err := json.Unmarshal(body, &u); err != nil {
		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	user, err := s.conn.SaveUser(ctx, u.ID, u.Username)
	if err != nil {
		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	uuid, _ := user.ID.Value()

	seshUser := SessionUser{
		ID:          uuid.(string),
		DiscordName: user.DiscordName,
		CreatedAt:   user.CreatedAt.Time,
		Roles:       user.Roles,
	}

	sesh.Set("user", seshUser)
	if err := sesh.Save(); err != nil {
		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	redirectLoggedIn(c)
}

func (s *Server) meHandler(c *gin.Context) {
	sesh := sessions.Default(c)

	u := ReturnedUser{}
	if val := sesh.Get("user"); val != nil {
		user := val.(*SessionUser)
		u.ID = user.ID
		u.Username = user.DiscordName
		u.CreatedAt = user.CreatedAt
	}
	fmt.Println(u)

	c.JSON(http.StatusOK, u)
}

func (s *Server) logout(c *gin.Context) {
	sesh := sessions.Default(c)
	sesh.Clear()
	if err := sesh.Save(); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	redirectLoggedIn(c)
}
