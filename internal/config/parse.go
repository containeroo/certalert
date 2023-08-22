package config

import (
	"certalert/internal/certificates"
	"certalert/internal/resolve"
	"certalert/internal/utils"
	"fmt"
	"net/url"
	"path/filepath"
	"reflect"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
)

// validateAuthConfig validates the auth config
func validateAuthConfig(authConfig Auth) error {
	basicValue := reflect.ValueOf(authConfig.Basic)
	bearerValue := reflect.ValueOf(authConfig.Bearer)

	basicZero := reflect.Zero(basicValue.Type())
	bearerZero := reflect.Zero(bearerValue.Type())

	if basicValue.Interface() != basicZero.Interface() && bearerValue.Interface() != bearerZero.Interface() {
		return fmt.Errorf("Both 'auth.basic' and 'auth.bearer' are defined")
	}

	return nil
}

// ParseConfig parse the config file and resolves variables
func ParseConfig(config *Config, failOnError bool) (err error) {
	// handleCertError is a helper function to handle errors during certificate validation
	handleCertError := func(cert certificates.Certificate, idx int, errMsg string) error {
		if failOnError {
			config.Certs[idx] = cert
			return fmt.Errorf(errMsg)
		}
		log.Warn(errMsg)
		return nil
	}

	// handlePushgatewayError is a helper function to handle errors during pushgateway validation
	handlePushgatewayError := func(errMsg string) error {
		if failOnError {
			return fmt.Errorf(errMsg)
		}
		log.Warn(errMsg)
		return nil
	}

	resolvedAddress, err := resolve.ResolveVariable(config.Pushgateway.Address)
	if err != nil {
		if err := handlePushgatewayError(fmt.Sprintf("Failed to resolve address for pushgateway: %v", err)); err != nil {
			return err
		}
	}
	if _, err := url.Parse(resolvedAddress); err != nil {
		if err := handlePushgatewayError(fmt.Sprintf("Invalid pushgateway address '%s': %v", resolvedAddress, err)); err != nil {
			return err
		}
	}
	config.Pushgateway.Address = resolvedAddress

	if err := validateAuthConfig(config.Pushgateway.Auth); err != nil {
		if err := handlePushgatewayError(err.Error()); err != nil {
			return err
		}
	}

	config.Pushgateway.Auth.Basic.Password, err = resolve.ResolveVariable(config.Pushgateway.Auth.Basic.Password)
	if err != nil {
		if err := handlePushgatewayError(fmt.Sprintf("Failed to resolve password for pushgateway: %v", err)); err != nil {
			return err
		}
	}

	config.Pushgateway.Auth.Bearer.Token, err = resolve.ResolveVariable(config.Pushgateway.Auth.Bearer.Token)
	if err != nil {
		if err := handlePushgatewayError(fmt.Sprintf("Failed to resolve token for pushgateway: %v", err)); err != nil {
			return err
		}
	}

	if config.Pushgateway.Job == "" {
		config.Pushgateway.Job = "certalert"
	} else {
		config.Pushgateway.Job, err = resolve.ResolveVariable(config.Pushgateway.Job)
		if err != nil {
			if err := handlePushgatewayError(fmt.Sprintf("Failed to resolve job for pushgateway: %v", err)); err != nil {
				return err
			}
		}
	}

	validFileExtensionTypes := utils.ExtractMapKeys(certificates.FileExtensionsToType)
	sort.Strings(validFileExtensionTypes) // sort the list to have a deterministic order

	for idx, cert := range config.Certs {
		if cert.Enabled != nil && !*cert.Enabled {
			log.Debugf("Skip certificate '%s' because is disabled", cert.Name)
			continue
		}

		if cert.Path == "" {
			if err := handleCertError(cert, idx, fmt.Sprintf("Certificate '%s' has no 'path' defined", cert.Name)); err != nil {
				return err
			}
		}

		if err := utils.CheckFileAccessibility(cert.Path); err != nil {
			if err := handleCertError(cert, idx, fmt.Sprintf("Certificate '%s' is not accessible: %v", cert.Name, err)); err != nil {
				return err
			}
		}

		if cert.Name == "" {
			file := filepath.Base(cert.Path)
			// replace dots, spaces and underscores with dashes
			cert.Name = strings.Map(func(r rune) rune {
				if r == '.' || r == ' ' || r == '_' {
					return '-'
				}
				return r
			}, file)
		}

		if cert.Type == "" {
			ext := filepath.Ext(cert.Path)     // try to guess the type based on the file extension
			ext = strings.TrimPrefix(ext, ".") // remove the dot

			if inferredType, ok := certificates.FileExtensionsToType[ext]; ok {
				cert.Type = inferredType
			} else {
				reason := "missing file extension."
				if ext != "" {
					reason = fmt.Sprintf("unclear file extension (.%s).", ext)
				}
				errMsg := fmt.Sprintf("Certificate '%s' has no 'type' defined. Type can't be inferred due to the %s", cert.Name, reason)
				return handleCertError(cert, idx, errMsg)
			}
		}

		if !utils.IsInList(cert.Type, validFileExtensionTypes) {
			return handleCertError(cert, idx, fmt.Sprintf("Certificate '%s' has an invalid 'type'. Must be one of: %s", cert.Name, strings.Join(validFileExtensionTypes, ", ")))
		}

		pw, err := resolve.ResolveVariable(cert.Password)
		if err != nil {
			if err := handleCertError(cert, idx, fmt.Sprintf("Certifacate '%s' has a non resolvable 'password'. %v", cert.Name, err)); err != nil {
				return err
			}
		}
		cert.Password = pw

		config.Certs[idx] = cert

	}
	return nil
}
