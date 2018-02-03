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

var createCmd = &cobra.Command{
	Use:   "create",
	Short: "creates a check",
	Long:  "This subcommand lets a user create a HTTP check",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if createURL == "" {
			return errors.New("specify a url to monitor")
		}

		if createWebhookURL == "" {
			return errors.New("specify a webhook url")
		}

		_, err := url.Parse(createURL)
		if err != nil {
			return err
		}

		_, err = url.Parse(createWebhookURL)
		if err != nil {
			return err
		}

		if createPlaybookURL != "" {
			_, err = url.Parse(createPlaybookURL)
			if err != nil {
				return err
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		surl, err := cmd.Root().PersistentFlags().GetString("server")
		if err != nil {
			log.Fatal(err)
		}

		checks := lokahi.NewChecksProtobufClient(surl, &http.Client{})

		chk, err := checks.Create(context.Background(), &lokahi.CreateOpts{
			Url:         createURL,
			WebhookUrl:  createWebhookURL,
			Every:       int32(createEvery),
			PlaybookUrl: createPlaybookURL,
		})
		if err != nil {
			log.Fatal(err)
		}

		data, err := json.MarshalIndent(chk, "", "  ")
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(string(data))
	},
}

var (
	createURL, createWebhookURL, createPlaybookURL string
	createEvery                                    int
)

func init() {
	rootCmd.AddCommand(createCmd)

	createCmd.Flags().StringVarP(&createURL, "url", "u", "", "URL to monitor")
	createCmd.Flags().StringVarP(&createWebhookURL, "webhook-url", "w", "", "webhook URL to post updates to")
	createCmd.Flags().IntVarP(&createEvery, "every", "e", 0, "")
	createCmd.Flags().StringVarP(&createPlaybookURL, "playbook-url", "p", "", "playbook URL with operational instructions")
}
