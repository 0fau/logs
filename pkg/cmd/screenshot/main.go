package main

import (
	"context"
	"github.com/0fau/logs/pkg/s3"
	"github.com/0fau/logs/pkg/screenshot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func config() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LBF")

	viper.MustBindEnv(
		"ENVIRONMENT",
	)

	viper.MustBindEnv(
		"FRONTEND_URL",

		"DATABASE_URL",

		"S3_ENDPOINT",
		"S3_BUCKET",
		"S3_ACCESS_KEY_ID",
		"S3_SECRET_ACCESS_KEY",
	)
}

func main() {
	cmd := &cobra.Command{
		Use: "logs-screenshot",
		RunE: func(cmd *cobra.Command, args []string) error {
			config()

			a := screenshot.NewServer(&screenshot.Config{
				FrontendURL: viper.GetString("FRONTEND_URL"),
				DatabaseURL: viper.GetString("DATABASE_URL"),
				S3: &s3.Config{
					Endpoint:        viper.GetString("S3_ENDPOINT"),
					Bucket:          viper.GetString("S3_BUCKET"),
					AccessKeyID:     viper.GetString("S3_ACCESS_KEY_ID"),
					SecretAccessKey: viper.GetString("S3_SECRET_ACCESS_KEY"),
				},
			})
			return a.Run(context.Background())
		},
	}

	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
