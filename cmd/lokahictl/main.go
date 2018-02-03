package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "lokahictl",
		Short: "Control lokahi, a http healthchecking service",
		Long:  "See https://github.com/Xe/lokahi for more information",
	}
)

func init() {
	rootCmd.PersistentFlags().String("server", "http://AzureDiamond:hunter2@127.0.0.1:24253", "http url of the lokahid instance")
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
