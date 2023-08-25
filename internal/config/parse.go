package config

import (
	"certalert/internal/certificates"
	"certalert/internal/resolve"
	"certalert/internal/utils"
	"fmt"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

// Parse parse the config file and resolves variables
func (c *Config) Parse() (err error) {
	if err := c.parsePushgatewayConfig(); err != nil {
		return err
	}

	if err := c.parseCertificatesConfig(); err != nil {
		return err
	}

	return nil
}

// parseCertificatesConfig validates the certificates config
func (c *Config) parseCertificatesConfig() (err error) {
	// handleError is a helper function to handle errors during certificate validation
	handleError := func(cert certificates.Certificate, idx int, errMsg string) error {
		if c.FailOnError {
			c.Certs[idx] = cert
			return fmt.Errorf(errMsg)
		}
		log.Warn(errMsg)
		return nil
	}

	for idx, cert := range c.Certs {
		if cert.Enabled != nil && !*cert.Enabled {
			log.Debugf("Skip certificate '%s' because is disabled", cert.Name)
			continue
		}

		if cert.Path == "" {
			if err := handleError(cert, idx, fmt.Sprintf("Certificate '%s' has no 'path' defined.", cert.Name)); err != nil {
				return err
			}
		}

		if err := utils.CheckFileAccessibility(cert.Path); err != nil {
			if err := handleError(cert, idx, fmt.Sprintf("Certificate '%s' is not accessible. %v", cert.Name, err)); err != nil {
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
				return handleError(cert, idx, errMsg)
			}
		}

		if !utils.IsInList(cert.Type, certificates.FileExtensionsTypes) {

			if err := handleError(cert, idx, fmt.Sprintf("Certificate '%s' has an invalid type '%s'. Must be one of '%s' or '%s'.", cert.Name, cert.Type, strings.Join(certificates.FileExtensionsTypes[:certificates.LenFileExtensionsTypes-1], "', '"), certificates.FileExtensionsTypes[certificates.LenFileExtensionsTypes-1])); err != nil {
				return err
			}
		}

		pw, err := resolve.ResolveVariable(cert.Password)
		if err != nil {
			if err := handleError(cert, idx, fmt.Sprintf("Certifacate '%s' has a non resolvable 'password'. %v", cert.Name, err)); err != nil {
				return err
			}
		}
		cert.Password = pw

		c.Certs[idx] = cert

	}
	return nil

}

// parsePushgatewayConfig validates the pushgateway config
func (c *Config) parsePushgatewayConfig() (err error) {
	// handlePushgatewayError is a helper function to handle errors during pushgateway validation
	handlePushgatewayError := func(errMsg string) error {
		if c.FailOnError {
			return fmt.Errorf(errMsg)
		}
		log.Warn(errMsg)
		return nil
	}
	resolvedAddress, err := resolve.ResolveVariable(c.Pushgateway.Address)
	if err != nil {
		if err := handlePushgatewayError(fmt.Sprintf("Failed to resolve address for pushgateway. %v", err)); err != nil {
			return err
		}
	}
	if resolvedAddress == "" && c.Pushgateway.Address != "" {
		if err := handlePushgatewayError("Pushgateway address was resolved to empty."); err != nil {
			return err
		}
	}
	if resolvedAddress != "" && !isValidURL(resolvedAddress) {
		if err := handlePushgatewayError(fmt.Sprintf("Invalid pushgateway address '%s'.", resolvedAddress)); err != nil {
			return err
		}
	}
	c.Pushgateway.Address = resolvedAddress

	if err := validateAuthConfig(c.Pushgateway.Auth); err != nil {
		if err := handlePushgatewayError(err.Error()); err != nil {
			return err
		}
	}

	c.Pushgateway.Auth.Basic.Password, err = resolve.ResolveVariable(c.Pushgateway.Auth.Basic.Password)
	if err != nil {
		if err := handlePushgatewayError(fmt.Sprintf("Failed to resolve basic auth password for pushgateway. %v", err)); err != nil {
			return err
		}
	}

	c.Pushgateway.Auth.Bearer.Token, err = resolve.ResolveVariable(c.Pushgateway.Auth.Bearer.Token)
	if err != nil {
		if err := handlePushgatewayError(fmt.Sprintf("Failed to resolve bearer token for pushgateway. %v", err)); err != nil {
			return err
		}
	}

	if c.Pushgateway.Job == "" {
		c.Pushgateway.Job = "certalert"
	} else {
		c.Pushgateway.Job, err = resolve.ResolveVariable(c.Pushgateway.Job)
		if err != nil {
			if err := handlePushgatewayError(fmt.Sprintf("Failed to resolve jobName for pushgateway. %v", err)); err != nil {
				return err
			}
		}
	}

	return nil
}
