package api

import (
	"context"
	"encoding/gob"
	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type ServerConfig struct {
	Address     string
	DatabaseURL string

	RedisAddress  string
	RedisPassword string
	SessionSecret string

	OAuth2 *oauth2.Config
}

type Server struct {
	conf   *ServerConfig
	router *gin.Engine

	queries *sql.Queries
}

func NewServer(conf *ServerConfig) *Server {
	router := gin.Default()
	router.MaxMultipartMemory = 8 << 20 // 8 MiB

	return &Server{
		conf:   conf,
		router: router,
	}
}

func (s *Server) Run(ctx context.Context) error {
	if err := database.Migrate(s.conf.DatabaseURL); err != nil {
		return err
	}

	pool, err := database.NewPool(ctx, s.conf.DatabaseURL)
	if err != nil {
		return err
	}
	defer pool.Close()
	s.queries = sql.New(pool)

	store, err := redis.NewStore(10, "tcp", s.conf.RedisAddress, s.conf.RedisPassword, []byte(s.conf.SessionSecret))
	if err != nil {
		return err
	}
	store.Options(sessions.Options{MaxAge: 604800}) // seven days
	s.router.Use(sessions.Sessions("sessions", store))

	gob.Register(sql.User{})
	s.router.POST("oauth2", s.oauth2)
	s.router.GET("oauth2/redirect", s.oauth2Redirect)
	s.router.GET("api/users/@me", s.meHandler)
	s.router.POST("logout", s.logout)

	s.router.POST("api/logs/upload", s.uploadHandler)

	return s.router.Run(s.conf.Address)
}
