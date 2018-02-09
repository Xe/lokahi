package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/Xe/uuid"
	"github.com/spf13/cobra"
)

var createLoadCmd = &cobra.Command{
	Use:   "create_load",
	Short: "creates a bunch of checks",
	Long:  "This subcommand lets a user create a HTTP check",
	Run: func(cmd *cobra.Command, args []string) {
		checks := lokahi.NewChecksProtobufClient(connServer, &http.Client{})

		for range make([]struct{}, 5000) {
			_, err := checks.Create(context.Background(), &lokahi.CreateOpts{
				Url:        "http://duke:9001?" + uuid.New(),
				WebhookUrl: "http://samplehook:9001/twirp/github.xe.lokahi.Webhook/Handle",
				Every:      60,
			})
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("5000 checks created")
	},
}

func init() {
	rootCmd.AddCommand(createLoadCmd)
}
