package api

import (
	"context"
	"encoding/gob"
	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/process"
	"github.com/0fau/logs/pkg/s3"
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

	S3     *s3.Config
	OAuth2 *oauth2.Config
}

type Server struct {
	config *ServerConfig
	router *gin.Engine

	processor *process.Processor
	conn      *database.DB
	s3        *s3.Client
}

func cors() gin.HandlerFunc {
	config := gincors.DefaultConfig()
	config.AllowOrigins = []string{"https://tauri.localhost", "http://localhost:5173"}
	config.AllowHeaders = []string{"access_token"}
	return gincors.New(config)
}

func NewServer(conf *ServerConfig) *Server {
	router := gin.Default()
	router.MaxMultipartMemory = 2 << 20 // 2 MiB
	router.Use(cors())

	return &Server{
		config: conf,
		router: router,
	}
}

func (s *Server) Run(ctx context.Context) error {
	conn, err := database.Connect(ctx, s.config.DatabaseURL, "logs", true)
	if err != nil {
		return errors.Wrap(err, "connecting to database")
	}
	s.conn = conn

	s.s3, err = s3.NewClient(s.config.S3)
	if err != nil {
		return errors.Wrap(err, "creating minio s3 client")
	}

	s.processor = process.NewLogProcessor(s.conn, s.s3)
	if err := s.processor.Initialize(); err != nil {
		return errors.Wrap(err, "initializing log processor")
	}

	store, err := redis.NewStore(10, "tcp", s.config.RedisAddress, s.config.RedisPassword, []byte(s.config.SessionSecret))
	if err != nil {
		return errors.Wrap(err, "creating redis sessions store")
	}
	store.Options(sessions.Options{MaxAge: 2628000}) // one month
	s.router.Use(sessions.Sessions("sessions", store))

	gob.Register(&SessionUser{})
	s.router.POST("oauth2", s.oauth2)
	s.router.GET("oauth2/redirect", s.oauth2Redirect)
	s.router.GET("api/users/@me", s.meHandler)
	s.router.POST("logout", s.logout)

	s.router.GET("images/avatar/:user", s.avatarHandler)

	s.router.GET("api/settings", s.settingsHandler)
	s.router.PUT("api/settings/username", s.setUsername)

	s.router.POST("api/logs", s.logs)
	s.router.POST("api/logs/upload", s.uploadHandler)
	// s.router.GET("api/logs/stats", s.statsHandler)
	s.router.GET("api/log/:log", s.logHandler)

	s.router.POST("api/users/@me/token", s.generateToken)

	return s.router.Run(s.config.Address)
}
