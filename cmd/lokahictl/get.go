package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/spf13/cobra"
)

var getCmd = &cobra.Command{
	Use:   "get <id>",
	Short: "dumps information about a check",
	Long:  "Gets information for a check by its unique ID. This will get as much information as possible.",

	PreRunE: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		surl, err := cmd.Root().PersistentFlags().GetString("server")
		if err != nil {
			return err
		}

		checks := lokahi.NewChecksProtobufClient(surl, &http.Client{})

		chk, err := checks.Get(ctx, &lokahi.CheckID{Id: args[0]})
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

func init() {
	rootCmd.AddCommand(getCmd)
}
