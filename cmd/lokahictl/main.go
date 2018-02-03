package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/Xe/lokahi/rpc/lokahi"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "lokahictl",
		Short: "Control lokahi, a http healthchecking service",
		Long:  "See https://github.com/Xe/lokahi for more information",
	}

	serverURL = rootCmd.Flags().String("server", "http://AzureDiamond:hunter2@127.0.0.1:24253", "http url of the lokahid instance")
	checks    lokahi.Checks
)

func init() {
	rootCmd.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		checks = lokahi.NewChecksProtobufClient(*serverURL, &http.Client{})
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
