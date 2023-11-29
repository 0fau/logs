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

func delete() *cobra.Command {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LCLI")

	viper.MustBindEnv("ADMIN_ADDRESS")
	addr := viper.GetString("ADMIN_ADDRESS")

	return &cobra.Command{
		Use:  "delete",
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
			_, err = cli.Delete(ctx, &admin.DeleteRequest{Encounter: int32(encID)})
			if err != nil {
				log.Fatal(errors.Wrap(err, "delete"))
			}
		},
	}
}
