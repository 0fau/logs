package api

import (
	"context"
	"encoding/gob"
	"github.com/0fau/logs/pkg/admin"
	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/process"
	"github.com/0fau/logs/pkg/s3"
	"github.com/cockroachdb/errors"
	gincors "github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"log"
)

type ServerConfig struct {
	Address     string
	DatabaseURL string

	RedisAddress  string
	RedisPassword string
	SessionSecret string

	S3     *s3.Config
	OAuth2 *oauth2.Config

	Admin *admin.Config
}

type Server struct {
	conf   *ServerConfig
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
		conf:   conf,
		router: router,
	}
}

func (s *Server) Run(ctx context.Context) error {
	conn, err := database.Connect(ctx, s.conf.DatabaseURL)
	if err != nil {
		return errors.Wrap(err, "connecting to database")
	}
	s.conn = conn

	s.s3, err = s3.NewClient(s.conf.S3)
	if err != nil {
		return errors.Wrap(err, "creating minio s3 client")
	}

	s.processor = process.NewLogProcessor(s.conn, s.s3)
	if err := s.processor.Initialize(); err != nil {
		return errors.Wrap(err, "initializing e processor")
	}

	a := admin.NewServer(s.conf.Admin, conn, s.s3, s.processor)
	go func() {
		if err := a.Run(); err != nil {
			log.Println(errors.Wrap(err, "running admin service"))
		}
	}()

	store, err := redis.NewStore(10, "tcp", s.conf.RedisAddress, s.conf.RedisPassword, []byte(s.conf.SessionSecret))
	if err != nil {
		return errors.Wrap(err, "creating redis sessions store")
	}
	store.Options(sessions.Options{MaxAge: 2628000}) // one month
	s.router.Use(sessions.Sessions("sessions", store))

	gob.Register(&SessionUser{})
	s.router.POST("oauth2", s.oauth2)
	s.router.GET("oauth2/redirect", s.oauth2Redirect)
	s.router.GET("api/users/@me", s.meHandler)
	s.router.GET("api/users/:user", s.userHandler)
	s.router.POST("logout", s.logout)

	s.router.GET("api/settings", s.settingsHandler)
	s.router.PUT("api/settings/username", s.setUsername)

	s.router.POST("api/logs", s.logs)
	s.router.POST("api/logs/upload", s.uploadHandler)
	// s.router.GET("api/logs/stats", s.statsHandler)
	s.router.GET("api/log/:log", s.logHandler)

	s.router.POST("api/users/@me/token", s.generateToken)

	return s.router.Run(s.conf.Address)
}
