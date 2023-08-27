package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/thanhpk/randstr"
	"io"
	"net/http"
)

type DiscordUser struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func redirectLoggedIn(c *gin.Context) {
	c.Redirect(http.StatusFound, "/me")
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
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if state.(string) != c.Query("state") {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	ctx := context.Background()
	token, err := s.conf.OAuth2.Exchange(ctx, c.Query("code"))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client := s.conf.OAuth2.Client(ctx, token)
	resp, err := client.Get("https://discord.com/api/users/@me")
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	u := DiscordUser{}
	if err := json.Unmarshal(body, &u); err != nil {
		fmt.Println("failed unmarhsal")
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	sesh.Set("user", u)
	if err := sesh.Save(); err != nil {
		fmt.Println("failed save", err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	redirectLoggedIn(c)
}

func (s *Server) meHandler(c *gin.Context) {
	sesh := sessions.Default(c)

	var u DiscordUser
	if val := sesh.Get("user"); val != nil {
		u = val.(DiscordUser)
	}

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
