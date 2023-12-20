/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>

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
	"certalert/internal/print"
	"certalert/internal/utils"
	"fmt"
	"os"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/spf13/cobra"
)

var (
	printAll               bool
	outputFormat           string
	supportedOutputFormats = []string{"text", "json", "yaml"}
)

// printCmd represents the print command
var printCmd = &cobra.Command{
	Use:   "print",
	Short: "Export certificates in different formats.",
	Long: `Prints certificates in different formats.

You can print all certificates or only a subset of them. The output format can be specified with the -o, --output flag.
The default output format is 'text'.

Examples:
	# Print all certificates in text format
	certalert print --all

	# Print all certificates in json format
	certalert print --all --output json

	# Print only the certificate with the name 'my-cert' in yaml format
	certalert print my-cert --output yaml
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if !utils.IsInList(outputFormat, supportedOutputFormats) {
			fmt.Printf("Unsupported output format: %s. Supported formats are: %s\n", outputFormat, strings.Join(supportedOutputFormats, ", "))
			cmd.Help()
			os.Exit(1)
		}

		// Parse config file in subcommand, because it is not needed for all subcommands
		// or there is a special order in which the flags should be parsed
		if err := config.App.Parse(); err != nil {
			log.Fatal().Msgf("Error parsing config file: %v", err)
		}

		if printAll {
			// Handle --all flag
			output, err := print.ConvertCertificatesToFormat(outputFormat, config.App.Certs, config.App.FailOnError)
			if err != nil {
				log.Fatal().Err(err)
			}
			fmt.Println(output)
			return
		}

		// Handle arguments
		if len(args) < 1 {
			fmt.Println("Please provide at least one argument or use the --all flag")
			cmd.Help()
			os.Exit(1)
		}

		var certs []certificates.Certificate
		// Create a list with all wanted certificates
		for _, arg := range args {
			certificate, err := certificates.GetCertificateByName(arg, config.App.Certs)
			if err != nil {
				log.Fatal().Err(err)
			}
			certs = append(certs, *certificate)
		}

		// Print the certificates
		output, err := print.ConvertCertificatesToFormat(outputFormat, certs, config.App.FailOnError)
		if err != nil {
			log.Fatal().Err(err)
		}
		fmt.Println(output)
	},
}

func init() {
	rootCmd.AddCommand(printCmd)

	printCmd.PersistentFlags().BoolVarP(&printAll, "all", "A", false, "Prints all certificates")

	printCmd.Flags().StringVarP(&outputFormat, "output", "o", "text", fmt.Sprintf("Output format. One of: %s", strings.Join(supportedOutputFormats, "|")))
	printCmd.RegisterFlagCompletionFunc("output", func(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		return supportedOutputFormats, cobra.ShellCompDirectiveDefault
	})
}
