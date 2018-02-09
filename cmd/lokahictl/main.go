package main

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	rootCmd = &cobra.Command{
		Use:   "lokahictl",
		Short: "Control lokahi, a http healthchecking service",
		Long:  "See https://github.com/Xe/lokahi for more information",
	}
)

var (
	cfgFile    string
	connServer string
)

func init() {
	cobra.OnInitialize(initConfig)
	rootCmd.PersistentFlags().StringVarP(&connServer, "server", "s", "http://127.0.0.1:24253", "http url of the lokahid instance")
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.lokahictl.hcl)")

	viper.BindPFlag("server", rootCmd.PersistentFlags().Lookup("server"))
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func initConfig() {
	// Don't forget to read config either from cfgFile or from home directory!
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".cobra" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".lokahictl")
	}

	// ignore errors
	_ = viper.ReadInConfig()

	if s := viper.GetString("server"); s != "" {
		connServer = s
	}
}
