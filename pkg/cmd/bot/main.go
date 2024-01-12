package main

import (
	"context"
	"github.com/0fau/logs/pkg/bot"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func config() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LBOT")

	viper.MustBindEnv(
		"ENVIRONMENT",
		"DATABASE_URL",
		"DISCORD_BOT_TOKEN",
		"DISCORD_GUILDID",
		"DISCORD_MESSAGEID",
		"DISCORD_ROLEID",
	)
}

func main() {
	cmd := &cobra.Command{
		Use: "logs-bot",
		RunE: func(cmd *cobra.Command, args []string) error {
			config()

			b := bot.Bot{
				DiscordConfig: bot.DiscordConfig{
					Token: viper.GetString("DISCORD_BOT_TOKEN"),

					GuildID:   viper.GetString("DISCORD_GUILDID"),
					MessageID: viper.GetString("DISCORD_MESSAGEID"),
					RoleID:    viper.GetString("DISCORD_ROLEID"),
				},
				DatabaseURL: viper.GetString("DATABASE_URL"),
			}
			return b.Run(context.Background())
		},
	}
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
