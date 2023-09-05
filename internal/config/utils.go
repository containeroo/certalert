package config

import "net/url"

// isValidURL tests a string to determine if it is a well-structured URL.
func isValidURL(str string) bool {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != ""
}
