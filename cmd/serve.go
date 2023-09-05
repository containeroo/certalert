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
	"certalert/internal/server"
	"certalert/internal/utils"

	"github.com/fsnotify/fsnotify"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listenAddress string

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch the web server to expose certificate metrics",
	Long: `The 'serve' command starts a web server that exposes certificate metrics. These metrics
are accessible at '/metrics' endpoint.

The web server's hostname and port can be defined using the --hostname and --port flags respectively.
The default hostname is 'localhost' and the default port is '8080'.

Example:
	certalert serve --hostname localhost --port 8080

The serve command also watches for changes in the application configuration file and reloads
the configuration if changes are detected.

Endpoints:
	- /: Home page
	- /-/reload: Reload the configuration file
	- /config: View the current configuration
	- /healthz: Health check endpoint
	- /metrics: Metrics endpoint

`,
	Run: func(cmd *cobra.Command, args []string) {

		// Watch for config changes
		viper.OnConfigChange(func(e fsnotify.Event) {
			log.Infof("Config file changed: %s", e.Name)

			if err := config.App.Read(viper.ConfigFileUsed()); err != nil {
				log.Fatalf("Error reading config file: %v", err)
			}

			if err := config.App.Parse(); err != nil {
				log.Fatalf("Unable to parse config: %s", err)
			}

			if err := utils.DeepCopy(config.App, &config.AppCopy); err != nil {
				log.Fatalf("Unable to copy config: %s", err)
			}

			if err := config.RedactConfig(&config.AppCopy); err != nil {
				log.Fatalf("Unable to redact config: %s", err)
			}

		})

		if config.App.AutoReloadConfig {
			log.Debug("Auto reloading configuration is enabled")
			viper.WatchConfig()
		}

		config.App.Version = version

		// this is only necessary if starting the web server
		if err := utils.DeepCopy(config.App, &config.AppCopy); err != nil {
			log.Fatalf("Unable to copy config: %s", err)
		}
		if err := config.RedactConfig(&config.AppCopy); err != nil {
			log.Fatalf("Unable to redact config: %s", err)
		}

		hostname, port, err := utils.ExtractHostAndPort(listenAddress)
		if err != nil {
			log.Fatalf("Unable to extract hostname and port: %s", err)
		}

		server.Run(hostname, port)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)

	serveCmd.Flags().StringVar(&listenAddress, "listen-address", ":8080", "The address to listen on for HTTP requests.")
	serveCmd.Flags().BoolVar(&config.App.AutoReloadConfig, "auto-reload-config", true, "Detects config changes and reloads the configuration file.")
}
