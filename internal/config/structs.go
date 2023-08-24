package config

import "certalert/internal/certificates"

// Config represents the config file
var App Config

// ConfigCopy represents the config file with sensitive data redacted
var AppCopy Config

var FailOnError bool

// Config represents the config file
type Config struct {
	AutoReloadConfig bool                       `mapstructure:"autoReloadConfig"`
	Version          string                     `mapstructure:"version"`
	Server           Server                     `mapstructure:"server"`
	Pushgateway      Pushgateway                `mapstructure:"pushgateway"`
	Certs            []certificates.Certificate `mapstructure:"certs"`
}

// Server represents the server config
type Server struct {
	Hostname string `mapstructure:"hostname"`
	Port     int    `mapstructure:"port"`
}

// Pushgateway represents the pushgateway config
type Pushgateway struct {
	Address            string `mapstructure:"address"`
	InsecureSkipVerify bool   `mapstructure:"insecureSkipVerify"`
	Job                string `mapstructure:"job"`
	Auth               Auth   `mapstructure:"auth,omitempty" yaml:"auth,omitempty"`
}

// Auth represents the pushgateway auth config
type Auth struct {
	Basic  Basic  `mapstructure:"basic,omitempty" yaml:"basic,omitempty"`
	Bearer Bearer `mapstructure:"bearer,omitempty" yaml:"bearer,omitempty"`
}

// Basic represents the pushgateway basic auth config
type Basic struct {
	Password string `mapstructure:"password"`
	Username string `mapstructure:"username"`
}

// Bearer represents the pushgateway bearer auth config
type Bearer struct {
	Token string `mapstructure:"token"`
}
