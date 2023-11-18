package main

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use: "logs",
	}
	cmd.AddCommand(extract())
	cmd.AddCommand(process())
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
