package main

import (
	"github.com/0fau/logs/pkg/cmd/cli/logs"
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LOGS_CLI")
	viper.MustBindEnv("ADMIN_ADDRESS")

	cmd := &cobra.Command{
		Use: "bm",
	}
	cmd.AddCommand(extract())
	cmd.AddCommand(role())
	cmd.AddCommand(process())
	cmd.AddCommand(processAll())
	cmd.AddCommand(delete())
	cmd.AddCommand(deleteUserLogs())
	cmd.AddCommand(runOperation())
	cmd.AddCommand(logs.Cmd())
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
