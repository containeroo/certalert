package config

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strconv"
)

// ExtractHostAndPort extracts the hostname and port from the listen address
func ExtractHostAndPort(address string) (string, int, error) {
	host, portStr, err := net.SplitHostPort(address)
	if err != nil {
		return "", 0, err
	}

	port, err := strconv.Atoi(portStr)
	if err != nil {
		return "", 0, err
	}

	return host, port, nil
}

// validateAuthConfig validates the auth config
func validateAuthConfig(authConfig Auth) error {
	basicValue := reflect.ValueOf(authConfig.Basic)
	bearerValue := reflect.ValueOf(authConfig.Bearer)

	basicZero := reflect.Zero(basicValue.Type())
	bearerZero := reflect.Zero(bearerValue.Type())

	if basicValue.Interface() != basicZero.Interface() && bearerValue.Interface() != bearerZero.Interface() {
		return fmt.Errorf("Both 'auth.basic' and 'auth.bearer' are defined.")
	}

	return nil
}

// isValidURL tests a string to determine if it is a well-structured URL.
func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
