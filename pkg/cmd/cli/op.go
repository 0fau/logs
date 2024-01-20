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
)

func runOperation() *cobra.Command {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LCLI")

	viper.MustBindEnv("ADMIN_ADDRESS")
	addr := viper.GetString("ADMIN_ADDRESS")

	return &cobra.Command{
		Use:  "run-op",
		Args: cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(errors.Wrap(err, "grpc dial"))
			}

			ctx := context.Background()
			cli := admin.NewAdminClient(conn)
			_, err = cli.RunOperation(ctx, &admin.RunOperationRequest{})
			if err != nil {
				log.Fatal(errors.Wrap(err, "process all"))
			}
		},
	}
}
