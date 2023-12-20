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
	"certalert/internal/certificates"
	"certalert/internal/config"
	"certalert/internal/pushgateway"
	"fmt"
	"os"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var pushAll bool

// pushCmd represents the push command
var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push certificate expiration as a epoch to a Prometheus Pushgateway instance",
	Long: `Push is a command that allows you to push the expiration as an epoch about certificates to a Prometheus Pushgateway instance.

The command can either push metadata for all certificates by using the --all flag or it can push metadata for specific certificates by specifying their names as command-line arguments.

If no arguments are provided and the --all flag is not set, the command will print a help message and exit.

Examples:
  # Push metadata for all certificates
  certalert push --all

  # Push metadata for a single certificate
  certalert push my-certificate-name

  # Push metadata for multiple certificates
  certalert push my-certificate-name another-certificate-name
`,

	Run: func(cmd *cobra.Command, args []string) {
		// Parse config file in subcommand, because it is not needed for all subcommands
		// or there is a special order in which the flags should be parsed
		if err := config.App.Parse(); err != nil {
			log.Fatal().Msgf("Error parsing config file: %v", err)
		}

		if pushAll {
			// Handle --all flag
			if err := pushgateway.Send(
				config.App.Pushgateway.Address,
				config.App.Pushgateway.Job,
				config.App.Pushgateway.Auth,
				config.App.Certs,
				config.App.Pushgateway.InsecureSkipVerify,
				config.App.FailOnError); err != nil {
				log.Fatal().Err(err)
			}
			return
		}

		// Handle arguments
		if len(args) < 1 {
			fmt.Println("Please provide at least one argument or use the --all flag")
			cmd.Help()
			os.Exit(1)
		}

		for _, arg := range args {
			certificate, err := certificates.GetCertificateByName(arg, config.App.Certs)
			if err != nil {
				log.Panic().Err(err)
			}
			if err := pushgateway.Send(
				config.App.Pushgateway.Address,
				config.App.Pushgateway.Job,
				config.App.Pushgateway.Auth,
				[]certificates.Certificate{*certificate},
				config.App.Pushgateway.InsecureSkipVerify,
				config.App.FailOnError); err != nil {
				log.Panic().Err(err)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(pushCmd)

	pushCmd.PersistentFlags().BoolVarP(&pushAll, "all", "A", false, "Push all certificates")
	pushCmd.PersistentFlags().BoolVarP(&config.App.Pushgateway.InsecureSkipVerify, "insecure-skip-verify", "i", false, "Skip TLS certificate verification")
}
