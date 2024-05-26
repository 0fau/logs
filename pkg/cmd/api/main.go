package main

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"golang.org/x/oauth2"

	"github.com/0fau/logs/pkg/api"
	"github.com/0fau/logs/pkg/s3"
)

func bindEnv() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LBF")

	viper.MustBindEnv(
		"ENVIRONMENT",
	)

	viper.MustBindEnv("API_SERVER_ADDRESS")
	viper.SetDefault("API_SERVER_ADDRESS", "0.0.0.0:3000")

	viper.MustBindEnv("API_SERVER_DATABASE_URL")

	viper.MustBindEnv("API_SERVER_REDIS_ADDRESS")
	viper.SetDefault("API_SERVER_REDIS_ADDRESS", "localhost:6379")

	viper.MustBindEnv("API_SERVER_REDIS_PASSWORD")
	viper.MustBindEnv("API_SERVER_SESSION_SECRET")

	viper.MustBindEnv(
		"API_SERVER_DISCORD_OAUTH2_CLIENT_ID",
		"API_SERVER_DISCORD_OAUTH2_CLIENT_SECRET",
	)

	viper.MustBindEnv(
		"S3_ENDPOINT",
		"S3_BUCKET",
		"S3_ACCESS_KEY_ID",
		"S3_SECRET_ACCESS_KEY",
	)

	viper.BindEnv(
		"PYROSCOPE_ENABLED",
		"PYROSCOPE_SERVER",
		"PYROSCOPE_USER",
		"PYROSCOPE_PASSWORD",
	)
}

func main() {
	cmd := &cobra.Command{
		Use: "logs-api",
		RunE: func(cmd *cobra.Command, args []string) error {
			bindEnv()

			s := api.NewServer(&api.ServerConfig{
				Address:       viper.GetString("API_SERVER_ADDRESS"),
				DatabaseURL:   viper.GetString("API_SERVER_DATABASE_URL"),
				SessionSecret: viper.GetString("API_SERVER_SESSION_SECRET"),
				RedisAddress:  viper.GetString("API_SERVER_REDIS_ADDRESS"),
				RedisPassword: viper.GetString("API_SERVER_REDIS_PASSWORD"),

				PyroscopeEnabled:  viper.GetBool("PYROSCOPE_ENABLED"),
				PyroscopeServer:   viper.GetString("PYROSCOPE_SERVER"),
				PyroscopeUser:     viper.GetString("PYROSCOPE_USER"),
				PyroscopePassword: viper.GetString("PYROSCOPE_PASSWORD"),

				OAuth2: &oauth2.Config{
					ClientID:     viper.GetString("API_SERVER_DISCORD_OAUTH2_CLIENT_ID"),
					ClientSecret: viper.GetString("API_SERVER_DISCORD_OAUTH2_CLIENT_SECRET"),
					RedirectURL:  viper.GetString("API_SERVER_DISCORD_OAUTH2_REDIRECT_URL"),
					Scopes:       []string{"identify"},

					Endpoint: oauth2.Endpoint{
						TokenURL: "https://discord.com/api/oauth2/token",
						AuthURL:  "https://discord.com/oauth2/authorize",
					},
				},

				S3: &s3.Config{
					Endpoint:        viper.GetString("S3_ENDPOINT"),
					Bucket:          viper.GetString("S3_BUCKET"),
					AccessKeyID:     viper.GetString("S3_ACCESS_KEY_ID"),
					SecretAccessKey: viper.GetString("S3_SECRET_ACCESS_KEY"),
				},
			})

			gin.SetMode(gin.ReleaseMode)

			return s.Run(context.Background())
		},
	}

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
