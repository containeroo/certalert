package config

import (
	"certalert/internal/certificates"
	"certalert/internal/resolve"
	"certalert/internal/utils"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// Parse parses the configuration settings from the specified sources.
// It calls helper methods to parse the Pushgateway and Certificates configurations.
// Additionally, it validates and extracts the hostname and port from the configured listen address.
//
// Returns:
//   - error
//     An error if there is any issue parsing or validating the configuration settings.
func (c *Config) Parse() (err error) {
	if err := c.parsePushgatewayConfig(); err != nil {
		return err
	}

	if err := c.parseCertificatesConfig(); err != nil {
		return err
	}

	_, _, err = utils.ExtractHostAndPort(c.Server.ListenAddress)
	if err != nil {
		return fmt.Errorf("Unable to extract hostname and port: %s", err)
	}

	return nil
}

// parseCertificatesConfig parses the Certificates configuration settings.
// It validates and processes each certificate entry in the configuration.
//
// Returns:
//   - error
//     An error if there is any issue parsing or validating the certificate configurations.
func (c *Config) parseCertificatesConfig() (err error) {
	// handleFailOnError is a helper function to handle errors during certificate validation
	handleFailOnError := func(cert certificates.Certificate, idx int, errMsg string) error {
		if c.FailOnError {
			c.Certs[idx] = cert
			return fmt.Errorf(errMsg)
		}
		log.Warn().Msg(errMsg)
		return nil
	}

	for idx, cert := range c.Certs {
		if cert.Enabled != nil && !*cert.Enabled {
			log.Debug().Msgf("Skip certificate '%s' because is disabled", cert.Name)
			continue
		}

		if cert.Path == "" {
			if err := handleFailOnError(cert, idx, fmt.Sprintf("Certificate '%s' has no 'path' defined.", cert.Name)); err != nil {
				return err
			}
		}

		if err := utils.CheckFileAccessibility(cert.Path); err != nil {
			if err := handleFailOnError(cert, idx, fmt.Sprintf("Certificate '%s' is not accessible. %v", cert.Name, err)); err != nil {
				return err
			}
		}

		if cert.Name == "" {
			file := filepath.Base(cert.Path)
			// Replace dots, spaces and underscores with dashes
			cert.Name = strings.Map(func(r rune) rune {
				if r == '.' || r == ' ' || r == '_' {
					return '-'
				}
				return r
			}, file)
		}

		if cert.Type == "" {
			ext := strings.TrimPrefix(filepath.Ext(cert.Path), ".") // extract file extation and remove leading dot
			if ext == "" {
				errMsg := fmt.Sprintf("Certificate '%s' has no 'type' defined and is missing a file extension.", cert.Name)
				return handleFailOnError(cert, idx, errMsg)
			}

			inferredType, ok := certificates.FileExtensionsToType[ext]
			if !ok {
				errMsg := fmt.Sprintf("Certificate '%s' has no 'type' defined. Type can't be inferred due to unclear file extension (.%s).", cert.Name, ext)
				return handleFailOnError(cert, idx, errMsg)
			}
			cert.Type = inferredType
		}

		// The Type can be specified in the config file, but it must be one of the supported types
		if !utils.IsInList(cert.Type, certificates.FileExtensionsTypesSorted) {
			if err := handleFailOnError(cert, idx, fmt.Sprintf("Certificate '%s' has an invalid type '%s'. Must be one of %s.", cert.Name, cert.Type, certificates.FileExtensionsTypesSorted)); err != nil {
				return err
			}
		}

		pw, err := resolve.ResolveVariable(cert.Password)
		if err != nil {
			if err := handleFailOnError(cert, idx, fmt.Sprintf("Certifacate '%s' has a non resolvable 'password'. %v", cert.Name, err)); err != nil {
				return err
			}
		}
		cert.Password = pw

		c.Certs[idx] = cert
	}

	return nil
}

// parsePushgatewayConfig parses the Pushgateway configuration settings.
// It validates and resolves variables in the pushgateway configuration.
//
// Returns:
//   - error
//     An error if there is any issue parsing or validating the pushgateway configuration.
func (c *Config) parsePushgatewayConfig() (err error) {
	// handleFailOnError is a helper function to handle errors during pushgateway validation
	handleFailOnError := func(errMsg string) error {
		if c.FailOnError {
			return fmt.Errorf(errMsg)
		}
		log.Warn().Msg(errMsg)
		return nil
	}
	if utils.HasStructField(c.Pushgateway, "Address") {
		resolvedAddress, err := resolve.ResolveVariable(c.Pushgateway.Address)
		if err != nil {
			if err := handleFailOnError(fmt.Sprintf("Failed to resolve address for pushgateway. %v", err)); err != nil {
				return err
			}
		}
		if resolvedAddress == "" && c.Pushgateway.Address != "" {
			if err := handleFailOnError("Pushgateway address was resolved to empty."); err != nil {
				return err
			}
		}
		if resolvedAddress != "" && !utils.IsValidURL(resolvedAddress) {
			if err := handleFailOnError(fmt.Sprintf("Invalid pushgateway address '%s'.", resolvedAddress)); err != nil {
				return err
			}
		}
		c.Pushgateway.Address = resolvedAddress
	}

	if err := c.Pushgateway.Auth.Validate(); err != nil {
		if err := handleFailOnError(err.Error()); err != nil {
			return err
		}
	}

	if utils.HasStructField(c.Pushgateway.Auth, "Basic.Username") {
		c.Pushgateway.Auth.Basic.Username, err = resolve.ResolveVariable(c.Pushgateway.Auth.Basic.Username)
		if err != nil {
			if err := handleFailOnError(fmt.Sprintf("Failed to resolve basic auth username for pushgateway. %v", err)); err != nil {
				return err
			}
		}
	}

	if utils.HasStructField(c.Pushgateway.Auth, "Basic.Password") {
		c.Pushgateway.Auth.Basic.Password, err = resolve.ResolveVariable(c.Pushgateway.Auth.Basic.Password)
		if err != nil {
			if err := handleFailOnError(fmt.Sprintf("Failed to resolve basic auth password for pushgateway. %v", err)); err != nil {
				return err
			}
		}
	}

	if utils.HasStructField(c.Pushgateway.Auth, "Bearer.Token") {
		c.Pushgateway.Auth.Bearer.Token, err = resolve.ResolveVariable(c.Pushgateway.Auth.Bearer.Token)
		if err != nil {
			if err := handleFailOnError(fmt.Sprintf("Failed to resolve bearer token for pushgateway. %v", err)); err != nil {
				return err
			}
		}
	}

	if utils.HasStructField(c.Pushgateway, "Job") {
		jobName := c.Pushgateway.Job
		if jobName == "" {
			jobName = "certalert"
		}

		c.Pushgateway.Job, err = resolve.ResolveVariable(jobName)
		if err != nil {
			if err := handleFailOnError(fmt.Sprintf("Failed to resolve jobName for pushgateway. %v", err)); err != nil {
				return err
			}
		}
	}

	return nil
}
