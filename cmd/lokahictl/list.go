package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/spf13/cobra"
)

var (
	listCount, listOffset int
	listStatus            bool
)

func init() {
	listCmd.Flags().IntVarP(&listCount, "count", "c", 30, "number of checks to return at once")
	listCmd.Flags().IntVarP(&listOffset, "offset", "o", 0, "number of checks to offset")
	listCmd.Flags().BoolVarP(&listStatus, "status", "s", false, "include detailed histogram status?")

	rootCmd.AddCommand(listCmd)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "lists all checks that you have permission to access",
	Long:  "Lists all information for all checks that you have permission to access",

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		checks := lokahi.NewChecksProtobufClient(connServer, &http.Client{})

		chk, err := checks.List(ctx, &lokahi.ListOpts{
			Count:         int32(listCount),
			Offset:        int32(listOffset),
			IncludeStatus: listStatus,
		})
		if err != nil {
			return err
		}

		data, err := json.MarshalIndent(chk, "", "  ")
		if err != nil {
			return err
		}

		fmt.Println(string(data))
		return nil
	},
}
