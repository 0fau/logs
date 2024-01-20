package api

import (
	"bytes"
	"context"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/0fau/logs/pkg/database/sql/structs"
	"github.com/bwmarrin/discordgo"
	crdbpgx "github.com/cockroachdb/cockroach-go/v2/crdb/crdbpgxv5"
	"github.com/cockroachdb/errors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
	"github.com/thanhpk/randstr"
	"image/png"
	"log"
	"net/http"
)

type SessionUser struct {
	ID         string
	DiscordTag string
	Avatar     bool
	Username   string
}

type ReturnedUser struct {
	ID         string `json:"id"`
	DiscordTag string `json:"discordTag"`
	Avatar     bool   `json:"avatar"`
	Username   string `json:"username"`
}

type ReturnedTokenUser struct {
	ID         string `json:"id"`
	DiscordTag string `json:"discordTag"`
	Avatar     bool   `json:"avatar"`
	Username   string `json:"username"`
	CanUpload  bool   `json:"canUpload"`
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
	url := s.config.OAuth2.AuthCodeURL(state)
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
	token, err := s.config.OAuth2.Exchange(ctx, c.Query("code"))
	if err != nil {
		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	client := s.config.OAuth2.Client(ctx, token)
	user, err := s.saveUser(ctx, client)
	if err != nil {
		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	uuid, _ := user.ID.Value()
	seshUser := SessionUser{
		ID:         uuid.(string),
		DiscordTag: user.DiscordTag,
		Avatar:     user.Avatar != "",
		Username:   user.Username.String,
	}

	sesh.Set("user", seshUser)
	if err := sesh.Save(); err != nil {
		log.Println(errors.WithStack(err))
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	redirectLoggedIn(c)
}

func (s *Server) saveUser(ctx context.Context, client *http.Client) (*sql.User, error) {
	dg, _ := discordgo.New("")
	dgUser, err := dg.User("@me", discordgo.WithClient(client))
	if err != nil {
		return nil, err
	}

	username := dgUser.Username
	if dgUser.Discriminator != "0" {
		username += "#" + dgUser.Discriminator
	}

	var user sql.User
	if err := crdbpgx.ExecuteTx(ctx, s.conn.Pool, pgx.TxOptions{}, func(tx pgx.Tx) error {
		qtx := s.conn.Queries.WithTx(tx)

		row, err := qtx.GetRolesByDiscordID(ctx, dgUser.ID)
		if err != nil && !errors.Is(err, pgx.ErrNoRows) {
			return errors.Wrap(err, "getting roles")
		}

		role, err := qtx.FetchWhitelist(ctx, dgUser.ID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return errors.Wrap(err, "fetching whitelist")
			}
		} else {
			row.Roles = append(row.Roles, role)
		}

		user, err = qtx.UpsertUser(ctx, sql.UpsertUserParams{
			DiscordID:  dgUser.ID,
			DiscordTag: username,
			Roles:      row.Roles,
			Avatar:     "",
			Settings: structs.UserSettings{
				SkipLanding:       false,
				LogVisibility:     "unlisted",
				ProfileVisibility: "unlisted",
			},
		})
		if err != nil {
			return errors.Wrap(err, "upserting user")
		}

		if user.Avatar != dgUser.Avatar {
			val, _ := user.ID.Value()
			uuid := val.(string)

			if dgUser.Avatar == "" {
				if err := s.s3.RemoveAvatar(ctx, uuid); err != nil {
					return errors.Wrap(err, "removing avatar from s3")
				}
			} else {
				img, err := dg.UserAvatarDecode(dgUser, discordgo.WithClient(client))
				if err != nil {
					return errors.Wrap(err, "user avatar decode")
				}

				buf := new(bytes.Buffer)
				if err := png.Encode(buf, img); err != nil {
					return errors.Wrap(err, "encoding avatar png")
				}

				if err := s.s3.SaveAvatar(ctx, uuid, buf.Bytes()); err != nil {
					return errors.Wrap(err, "saving avatar to s3")
				}
			}

			if err := qtx.UpdateAvatar(ctx, sql.UpdateAvatarParams{
				ID:     user.ID,
				Avatar: dgUser.Avatar,
			}); err != nil {
				return errors.Wrap(err, "updating avatar in db")
			}
		}

		return nil
	}); err != nil {
		return nil, errors.Wrap(err, "executing transaction")
	}

	return &user, nil
}

func (s *Server) meHandler(c *gin.Context) {
	token := c.GetHeader("access_token")
	if token != "" {
		user, err := s.conn.UserByAccessToken(context.Background(), token)
		if err != nil {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		id, _ := user.ID.Value()
		username := ""
		if name, err := user.Username.Value(); err == nil && name != nil {
			username = name.(string)
		}

		c.JSON(http.StatusOK, ReturnedTokenUser{
			ID:         id.(string),
			Username:   username,
			DiscordTag: user.DiscordTag,
			Avatar:     user.Avatar != "",
			CanUpload:  hasRoles(user.Roles, "alpha", "trusted"),
		})
		return
	}

	sesh := sessions.Default(c)
	if val := sesh.Get("user"); val != nil {
		user := val.(*SessionUser)
		c.JSON(http.StatusOK, ReturnedUser{
			ID:         user.ID,
			Username:   user.Username,
			DiscordTag: user.DiscordTag,
			Avatar:     user.Avatar,
		})
	} else {
		c.JSON(http.StatusUnauthorized, struct{}{})
	}
}

func (s *Server) logout(c *gin.Context) {
	sesh := sessions.Default(c)
	sesh.Clear()
	if err := sesh.Save(); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	redirectLoggedIn(c)
}
