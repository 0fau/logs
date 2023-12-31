package main

import (
	"context"
	"github.com/0fau/logs/pkg/admin"
	"github.com/cockroachdb/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"strconv"
)

func process() *cobra.Command {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LCLI")

	viper.MustBindEnv("ADMIN_ADDRESS")
	addr := viper.GetString("ADMIN_ADDRESS")

	return &cobra.Command{
		Use:  "process",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			encID, err := strconv.Atoi(args[0])
			if err != nil {
				log.Fatal(errors.Wrap(err, "convert encounter id"))
			}

			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(errors.Wrap(err, "grpc dial"))
			}

			ctx := context.Background()
			cli := admin.NewAdminClient(conn)
			_, err = cli.Process(ctx, &admin.ProcessRequest{Encounter: int32(encID)})
			if err != nil {
				log.Fatal(errors.Wrap(err, "process"))
			}
		},
	}
}

func processAll() *cobra.Command {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LCLI")

	viper.MustBindEnv("ADMIN_ADDRESS")
	addr := viper.GetString("ADMIN_ADDRESS")

	return &cobra.Command{
		Use:  "process-all",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(errors.Wrap(err, "grpc dial"))
			}

			ctx := context.Background()
			cli := admin.NewAdminClient(conn)
			_, err = cli.ProcessAll(ctx, &admin.ProcessAllRequest{})
			if err != nil {
				log.Fatal(errors.Wrap(err, "process all"))
			}
		},
	}
}
