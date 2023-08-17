package main

import (
	"context"
	"github.com/0fau/logs/pkg/api"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func config() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LBF")

	viper.MustBindEnv("API_SERVER_ADDRESS")
	viper.SetDefault("API_SERVER_ADDRESS", "0.0.0.0:3000")

	viper.MustBindEnv("API_SERVER_DATABASE_URL")

	viper.MustBindEnv("API_SERVER_REDIS_ADDRESS")
	viper.SetDefault("API_SERVER_REDIS_ADDRESS", "localhost:6379")

	viper.MustBindEnv("API_SERVER_REDIS_PASSWORD")
	viper.MustBindEnv("API_SERVER_SESSION_SECRET")
}

func main() {
	cmd := &cobra.Command{
		Use: "logs-api",
		RunE: func(cmd *cobra.Command, args []string) error {
			config()

			s := api.NewServer(&api.ServerConfig{
				Address:       viper.GetString("API_SERVER_ADDRESS"),
				DatabaseURL:   viper.GetString("API_SERVER_DATABASE_URL"),
				SessionSecret: viper.GetString("API_SERVER_SESSION_SECRET"),
				RedisAddress:  viper.GetString("API_SERVER_REDIS_ADDRESS"),
				RedisPassword: viper.GetString("API_SERVER_REDIS_PASSWORD"),
			})
			return s.Run(context.Background())
		},
	}

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
