package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Xe/lokahi/rpc/lokahiadmin"
	"github.com/spf13/cobra"
)

var runStatsCmd = &cobra.Command{
	Use:   "runstats",
	Short: "gets performance information",
	Long:  "dumps statistical performance information about all of the checks that have run since app boot",

	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		surl, err := cmd.Root().PersistentFlags().GetString("server")
		if err != nil {
			return err
		}

		rl := lokahiadmin.NewRunLocalProtobufClient(surl, &http.Client{})
		chk, err := rl.Stats(ctx, &lokahiadmin.Nil{})
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
	rootCmd.AddCommand(runStatsCmd)
}
