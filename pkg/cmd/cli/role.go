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

func role() *cobra.Command {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("LCLI")

	viper.MustBindEnv("ADMIN_ADDRESS")
	addr := viper.GetString("ADMIN_ADDRESS")

	return &cobra.Command{
		Use:  "role",
		Args: cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
			if err != nil {
				log.Fatal(errors.Wrap(err, "grpc dial"))
			}

			role := ""
			if args[0] == "trust" {
				role = "trusted"
			} else {
				log.Fatal("unknown role")
			}

			ctx := context.Background()
			cli := admin.NewAdminClient(conn)
			_, err = cli.Role(ctx, &admin.RoleRequest{
				Action:  admin.RoleRequest_Add,
				Role:    role,
				Discord: args[1],
			})
			if err != nil {
				log.Fatal(errors.Wrap(err, "process"))
			}
		},
	}
}
