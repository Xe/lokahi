package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/spf13/cobra"
)

var (
	runIDs string
)

var runCmd = &cobra.Command{
	Use:   "run [check-id, [check-id]]",
	Short: "runs a check",
	Long:  "run a check by id or list of ids",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		surl, err := cmd.Root().PersistentFlags().GetString("server")
		if err != nil {
			log.Fatal(err)
		}

		rl := lokahi.NewRunsProtobufClient(surl, &http.Client{})

		cids := &lokahi.CheckIDs{}

		for _, id := range args {
			cids.Ids = append(cids.Ids, id)
		}

		run, err := rl.Checks(ctx, cids)
		if err != nil {
			log.Fatal(err)
		}

		data, err := json.MarshalIndent(run, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(data))
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
