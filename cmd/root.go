/*
Copyright © 2023 containeroo hello©containeroo.ch

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"certalert/internal/config"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

const (
	version = "v0.0.11"
)

var cfgFile string
var verbose, silent, printVersion bool

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "certalert",
	Short: "CertAlert is a tool to monitor the expiration dates of digital certificates",
	Long: `CertAlert can handle a variety of certificate types, including .p12, .pkcs12, .pem, .crt, and .jks files.

	You can execute specific commands for different actions:
	1. Use the 'push' command to manually push metrics to the Prometheus Pushgateway.
	2. Use the 'serve' command to start a server that provides a '/metrics' endpoint for Prometheus to scrape.

	For a full list of commands and options, use 'certalert --help'.
	`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Enter here before any subcommand is executed

		if verbose {
			log.SetLevel(log.DebugLevel)
			log.Debugf("Verbose output enabled")
		} else if silent {
			log.SetLevel(log.ErrorLevel)
			log.Debugf("Silent output enabled")
		}

		if printVersion {
			fmt.Println("CertAlert version:", version)
			os.Exit(0)
		}

		if err := config.ParseConfig(&config.App); err != nil {
			log.Fatalf("Error parsing config file: %v", err)
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Enter here if no subcommand is specified
		if printVersion {
			fmt.Println("CertAlert version:", version)
			os.Exit(0)
		} else {
			cmd.Help()
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $HOME/.certalert.yaml)")

	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "verbose output")
	rootCmd.PersistentFlags().BoolVarP(&silent, "silent", "s", false, "silent output")
	rootCmd.MarkFlagsMutuallyExclusive("verbose", "silent")

	rootCmd.PersistentFlags().BoolVarP(&config.FailOnError, "fail-on-error", "f", false, "fail on error")
	rootCmd.PersistentFlags().BoolVarP(&printVersion, "version", "V", false, "print version and exit")

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".certalert" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".certalert")
	}

	viper.AutomaticEnv() // read in environment variables that match

	err := config.ReadConfigFile(viper.GetViper().ConfigFileUsed(), &config.App)
	if err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}
}
