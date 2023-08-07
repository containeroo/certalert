package config

import (
	"certalert/internal/certificates"
	"certalert/internal/utils"
	"fmt"
	"path/filepath"
	"reflect"
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
func ParseConfig(app *Config) (err error) {
	app.Pushgateway.Address, err = utils.ResolveVariable(app.Pushgateway.Address)
	if err != nil {
		return fmt.Errorf("Failed to resolve address for pushgateway: %v", err)
	}

	if err := validateAuthConfig(app.Pushgateway.Auth); err != nil {
		return err
	}

	app.Pushgateway.Auth.Basic.Password, err = utils.ResolveVariable(app.Pushgateway.Auth.Basic.Password)
	if err != nil {
		return fmt.Errorf("Failed to resolve password for pushgateway: %v", err)
	}

	app.Pushgateway.Auth.Bearer.Token, err = utils.ResolveVariable(app.Pushgateway.Auth.Bearer.Token)
	if err != nil {
		return fmt.Errorf("Failed to resolve token for pushgateway: %v", err)
	}

	if app.Pushgateway.Job == "" {
		app.Pushgateway.Job = "certalert"
	} else {
		app.Pushgateway.Job, err = utils.ResolveVariable(app.Pushgateway.Job)
		if err != nil {
			return fmt.Errorf("Failed to resolve job for pushgateway: %v", err)
		}
	}

	for idx, cert := range app.Certs {
		if cert.Enabled != nil && !*cert.Enabled {
			app.Certs[idx] = cert // update the certificate in the slice (maybe has changed from enabled to disabled)
			log.Debugf("Skip certificate '%s' because is disabled", cert.Name)
			continue
		}

		if cert.Path == "" {
			return fmt.Errorf("Certificate '%s' has no 'path' defined", cert.Name)
		}

		if err := utils.CheckFileAccessibility(cert.Path); err != nil {
			return fmt.Errorf("Certificate '%s' has an invalid 'path'. %s", cert.Name, err)
		}

		if cert.Name == "" {
			file := filepath.Base(cert.Path)
			name := strings.ReplaceAll(file, ".", "-")
			name = strings.ReplaceAll(file, " ", "-")
			name = strings.ReplaceAll(file, "_", "-")
			cert.Name = name
		}

		if cert.Type == "" {
			// try to guess the type based on the file extension
			ext := filepath.Ext(cert.Path)
			switch ext {
			case ".p12", ".pkcs12", ".pfx":
				cert.Type = "p12"
			case ".pem", ".crt":
				cert.Type = "pem"
			case ".jks":
				cert.Type = "jks"
			default:
				reason := "missing file extension."
				if ext != "" {
					reason = fmt.Sprintf("unclear file extension (%s).", ext)
				}

				return fmt.Errorf("Certificate '%s' has no 'type' defined. Type can't be inferred due to the %s", cert.Name, reason)
			}
		}
		if !utils.IsInList(cert.Type, certificates.ValidTypes) {
			return fmt.Errorf("Certificate '%s' has an invalid 'type'. Must be one of: %s", cert.Name, strings.Join(certificates.ValidTypes, ", "))
		}

		pw, err := utils.ResolveVariable(cert.Password)
		if err != nil {
			return fmt.Errorf("Certifacate '%s' cannot resolve 'password'. %v", cert.Name, err)
		}
		cert.Password = pw

		app.Certs[idx] = cert
	}

	return nil
}
