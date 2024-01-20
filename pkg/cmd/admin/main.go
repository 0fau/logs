package main

import (
	"context"
	"github.com/0fau/logs/pkg/admin"
	"github.com/0fau/logs/pkg/s3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func config() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LBF")

	viper.MustBindEnv(
		"ENVIRONMENT",
	)

	viper.MustBindEnv("API_SERVER_DATABASE_URL")
	viper.MustBindEnv(
		"S3_ENDPOINT",
		"S3_BUCKET",
		"S3_ACCESS_KEY_ID",
		"S3_SECRET_ACCESS_KEY",
	)

	viper.MustBindEnv("DISCORD_BOT_TOKEN")

	viper.MustBindEnv(
		"ADMIN_ADDRESS",
	)
}

func main() {
	cmd := &cobra.Command{
		Use: "logs-admin",
		RunE: func(cmd *cobra.Command, args []string) error {
			config()

			a := admin.NewServer(&admin.Config{
				Address:     viper.GetString("ADMIN_ADDRESS"),
				DatabaseURL: viper.GetString("API_SERVER_DATABASE_URL"),
				S3: &s3.Config{
					Endpoint:        viper.GetString("S3_ENDPOINT"),
					Bucket:          viper.GetString("S3_BUCKET"),
					AccessKeyID:     viper.GetString("S3_ACCESS_KEY_ID"),
					SecretAccessKey: viper.GetString("S3_SECRET_ACCESS_KEY"),
				},
				DiscordBotToken: viper.GetString("DISCORD_BOT_TOKEN"),
			})
			return a.Run(context.Background())
		},
	}

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
