package api

import (
	"context"
	"encoding/gob"
	"log"
	"runtime"

	"github.com/cockroachdb/errors"
	gincors "github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	ginredis "github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"github.com/grafana/pyroscope-go"
	"github.com/redis/go-redis/v9"
	"golang.org/x/oauth2"

	"github.com/0fau/logs/pkg/database"
	"github.com/0fau/logs/pkg/process"
	"github.com/0fau/logs/pkg/s3"
)

type ServerConfig struct {
	Address     string
	DatabaseURL string

	RedisAddress  string
	RedisPassword string
	SessionSecret string

	S3     *s3.Config
	OAuth2 *oauth2.Config

	PyroscopeEnabled  bool
	PyroscopeServer   string
	PyroscopeUser     string
	PyroscopePassword string
}

type Server struct {
	config *ServerConfig
	router *gin.Engine

	processor *process.Processor
	conn      *database.DB
	redis     *redis.Client
	s3        *s3.Client

	//tokens      map[string]APIToken
	//tokensMutex sync.RWMutex
}

func cors() gin.HandlerFunc {
	config := gincors.DefaultConfig()
	config.AllowOrigins = []string{
		"https://tauri.localhost",
		"http://localhost:5173",
		"http://localhost:5174",
	}
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
	if s.config.PyroscopeEnabled {
		go s.StartPyroscope()
	}

	conn, err := database.Connect(ctx, s.config.DatabaseURL, "logs", true)
	if err != nil {
		return errors.Wrap(err, "connecting to database")
	}
	s.conn = conn

	s.redis = redis.NewClient(&redis.Options{
		Addr:     s.config.RedisAddress,
		Password: s.config.RedisPassword,
		DB:       5,
	})

	s.s3, err = s3.NewClient(s.config.S3)
	if err != nil {
		return errors.Wrap(err, "creating minio s3 client")
	}

	s.processor = process.NewLogProcessor(s.conn, s.s3)
	if err := s.processor.Initialize(); err != nil {
		return errors.Wrap(err, "initializing log processor")
	}

	store, err := ginredis.NewStore(10, "tcp", s.config.RedisAddress, s.config.RedisPassword, []byte(s.config.SessionSecret))
	if err != nil {
		return errors.Wrap(err, "creating redis sessions store")
	}
	store.Options(sessions.Options{Path: "/", MaxAge: 2628000}) // one month
	s.router.Use(sessions.Sessions("session", store))

	gob.Register(&SessionUser{})
	s.router.POST("oauth2", s.oauth2)
	s.router.GET("oauth2/redirect", s.oauth2Redirect)
	s.router.GET("api/users/@me", s.meHandler)
	s.router.POST("logout", s.logout)

	images := s.router.Group("images")
	images.GET("avatar/:user", s.avatarHandler)
	images.GET("thumbnail/:log", s.thumbnailHandler)

	s.router.GET("api/settings", s.settingsHandler)
	s.router.PATCH("api/settings", s.updateSettings)
	s.router.PUT("api/settings/username", s.setUsername)

	s.router.POST("api/logs", s.logs)
	s.router.POST("api/logs/upload", s.uploadHandler)
	s.router.GET("api/log/:log", s.logHandler)
	s.router.GET("api/log/:log/short", s.shortLogHandler)
	s.router.PATCH("api/log/:log/settings", s.updateLogSettings)

	s.router.GET("api/profile", s.profileHandler)
	s.router.GET("api/profile/:username", s.userProfileHandler)

	s.router.POST("api/users/@me/token", s.generateToken)

	return s.router.Run(s.config.Address)
}

func (s *Server) StartPyroscope() {
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)

	if _, err := pyroscope.Start(pyroscope.Config{
		ApplicationName: "logs.fau.dev",

		Logger: nil,

		ServerAddress:     s.config.PyroscopeServer,
		BasicAuthUser:     s.config.PyroscopeUser,
		BasicAuthPassword: s.config.PyroscopePassword,

		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,

			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	}); err != nil {
		log.Fatal(err)
	}
}
