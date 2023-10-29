package api

import (
	"context"
	"encoding/gob"
	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/process"
	"github.com/cockroachdb/errors"
	gincors "github.com/gin-contrib/cors"
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

	processor *process.Processor
	conn      *database.DB
}

func cors() gin.HandlerFunc {
	config := gincors.DefaultConfig()
	config.AllowOrigins = []string{"https://tauri.localhost"}
	config.AllowHeaders = []string{"access_token"}
	return gincors.New(config)
}

func NewServer(conf *ServerConfig) *Server {
	router := gin.Default()
	router.MaxMultipartMemory = 2 << 20 // 2 MiB
	router.Use(cors())

	return &Server{
		conf:   conf,
		router: router,
	}
}

func (s *Server) Run(ctx context.Context) error {
	s.processor = process.NewLogProcessor()
	if err := s.processor.Initialize(); err != nil {
		return errors.Wrap(err, "initializing e processor")
	}

	conn, err := database.Connect(ctx, s.conf.DatabaseURL)
	if err != nil {
		return errors.Wrap(err, "connecting to database")
	}
	s.conn = conn

	store, err := redis.NewStore(10, "tcp", s.conf.RedisAddress, s.conf.RedisPassword, []byte(s.conf.SessionSecret))
	if err != nil {
		return errors.Wrap(err, "creating redis sessions store")
	}
	store.Options(sessions.Options{MaxAge: 604800}) // seven days
	s.router.Use(sessions.Sessions("sessions", store))

	gob.Register(&SessionUser{})
	s.router.POST("oauth2", s.oauth2)
	s.router.GET("oauth2/redirect", s.oauth2Redirect)
	s.router.GET("api/users/@me", s.meHandler)
	s.router.GET("api/users/:user", s.userHandler)
	s.router.POST("logout", s.logout)

	s.router.GET("api/settings", s.settingsHandler)
	s.router.PUT("api/settings/username", s.setUsername)

	s.router.POST("api/logs/upload", s.uploadHandler)
	s.router.GET("api/logs/@recent", s.recentLogs)
	//s.router.GET("api/logs/:log", s.logHandler)
	s.router.POST("api/users/@me/token", s.generateToken)

	return s.router.Run(s.conf.Address)
}
