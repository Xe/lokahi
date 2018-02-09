package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"

	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/spf13/cobra"
)

var putCmd = &cobra.Command{
	Use:   "put",
	Short: "puts updates to a check",
	Long:  "This subcommand lets a user put a HTTP check",
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if putURL == "" {
			return errors.New("specify a url to monitor")
		}

		if putWebhookURL == "" {
			return errors.New("specify a webhook url")
		}

		_, err := url.Parse(putURL)
		if err != nil {
			return err
		}

		_, err = url.Parse(putWebhookURL)
		if err != nil {
			return err
		}

		if putPlaybookURL != "" {
			_, err = url.Parse(putPlaybookURL)
			if err != nil {
				return err
			}
		}

		return cobra.ExactArgs(1)(cmd, args)
	},
	RunE: func(cmd *cobra.Command, args []string) error {
		checks := lokahi.NewChecksProtobufClient(connServer, &http.Client{})

		chk, err := checks.Get(context.Background(), &lokahi.CheckID{Id: args[0]})
		if err != nil {
			return err
		}

		if putURL != "" {
			chk.Url = putURL
		}

		if putWebhookURL != "" {
			chk.WebhookUrl = putWebhookURL
		}

		if putEvery != 0 {
			chk.Every = int32(putEvery)
		}

		if putPlaybookURL != "" {
			chk.PlaybookUrl = putPlaybookURL
		}

		_, err = checks.Put(context.Background(), chk)
		if err != nil {
			log.Fatal(err)
		}

		data, err := json.MarshalIndent(chk, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(data))
		return nil
	},
}

var (
	putURL, putWebhookURL, putPlaybookURL string
	putEvery                              int
)

func init() {
	putCmd.Flags().StringVarP(&putURL, "url", "u", "", "URL to monitor")
	putCmd.Flags().StringVarP(&putWebhookURL, "webhook-url", "w", "", "webhook URL to post updates to")
	putCmd.Flags().IntVarP(&putEvery, "every", "e", 0, "")
	putCmd.Flags().StringVarP(&putPlaybookURL, "playbook-url", "p", "", "playbook URL with operational instructions")

	rootCmd.AddCommand(putCmd)
}
