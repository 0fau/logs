package api

import (
	"context"
	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/database/sql"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
)

type ServerConfig struct {
	Address     string
	DatabaseURL string

	RedisAddress  string
	RedisPassword string
	SessionSecret string
}

type Server struct {
	conf   *ServerConfig
	router *gin.Engine

	queries *sql.Queries
}

func NewServer(conf *ServerConfig) *Server {
	router := gin.Default()

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
	s.router.Use(sessions.Sessions("sessions", store))

	return s.router.Run(s.conf.Address)
}
