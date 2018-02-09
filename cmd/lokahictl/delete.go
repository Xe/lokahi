package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/spf13/cobra"
)

var (
	deleteIDs string
)

var deleteCmd = &cobra.Command{
	Use:   "delete [check-id, [check-id]]",
	Short: "deletes a check",
	Long:  "delete a check by id or list of ids",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		checks := lokahi.NewChecksProtobufClient(connServer, &http.Client{})
		fails := map[string]string{}

		for _, id := range args {
			_, err := checks.Delete(ctx, &lokahi.CheckID{Id: id})
			if err != nil {
				fails[id] = err.Error()
			}
		}

		if len(fails) != 0 {
			fmt.Println("one or more deletions failed")

			for k, v := range fails {
				fmt.Printf("%s: %s\n", k, v)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(deleteCmd)
}
