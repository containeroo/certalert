package config

import (
	"certalert/internal/certificates"
	"fmt"
)

// Config represents the config file
var App Config

// AppCopy represents the config file with sensitive data redacted
var AppCopy Config

// Config represents the config file
type Config struct {
	Version          string                     `mapstructure:"version"`
	AutoReloadConfig bool                       `mapstructure:"autoReloadConfig,omitempty" "yaml:"autoReloadConfig,omitempty"`
	FailOnError      bool                       `mapstructure:"failOnError,omitempty "yaml:"failOnError,omitempty"`
	Server           Server                     `mapstructure:"server,omitempty" yaml:"server,omitempty"`
	Pushgateway      Pushgateway                `mapstructure:"pushgateway,omitempty" yaml:"pushgateway,omitempty"`
	Certs            []certificates.Certificate `mapstructure:"certs"`
}

// Server represents the server config
type Server struct {
	Hostname string `mapstructure:"hostname,omitempty" yaml:"hostname,omitempty"`
	Port     int    `mapstructure:"port,omitempty" yaml:"port,omitempty"`
}

// Pushgateway represents the pushgateway config
type Pushgateway struct {
	Address            string `mapstructure:"address,omitempty" yaml:"address,omitempty"`
	InsecureSkipVerify bool   `mapstructure:"insecureSkipVerify,omitempty" yaml:"insecureSkipVerify,omitempty"`
	Job                string `mapstructure:"job,omitempty" yaml:"job,omitempty"`
	Auth               Auth   `mapstructure:"auth,omitempty" yaml:"auth,omitempty"`
}

// Auth represents the pushgateway auth config
type Auth struct {
	Basic  *Basic  `mapstructure:"basic,omitempty" yaml:"basic,omitempty"`
	Bearer *Bearer `mapstructure:"bearer,omitempty" yaml:"bearer,omitempty"`
}

// Validate checks if basic auth and bearer auth are both defined
func (a *Auth) Validate() error {
	if a.Basic != nil && a.Bearer != nil {
		return fmt.Errorf("Both 'auth.basic' and 'auth.bearer' are defined.")
	}

	return nil
}

// Basic represents the pushgateway basic auth config
type Basic struct {
	Password string `mapstructure:"password,omitempty" yaml:"password,omitempty"`
	Username string `mapstructure:"username,omitempty" yaml:"username,omitempty"`
}

// Bearer represents the pushgateway bearer auth config
type Bearer struct {
	Token string `mapstructure:"token,omitempty" yaml:"token,omitempty"`
}
