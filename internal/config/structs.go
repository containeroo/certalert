package config

import "certalert/internal/certificates"

// Config represents the config file
var App Config

// ConfigCopy represents the config file with sensitive data redacted
var AppCopy Config

var FailOnError bool

// Config represents the config file
type Config struct {
	Server      Server                     `json:"server"`
	Pushgateway Pushgateway                `json:"pushgateway"`
	Certs       []certificates.Certificate `json:"certs"`
}

// Server represents the server config
type Server struct {
	Hostname string `json:"hostname"`
	Port     int    `json:"port"`
}

// Pushgateway represents the pushgateway config
type Pushgateway struct {
	Address            string `json:"address"`
	InsecureSkipVerify bool   `json:"insecureSkipVerify"`
	Job                string `json:"job"`
	Auth               Auth   `json:"auth"`
}

// Auth represents the pushgateway auth config
type Auth struct {
	Basic  Basic  `json:"basic,omitempty" yaml:"basic,omitempty"`
	Bearer Bearer `json:"bearer,omitempty" yaml:"bearer,omitempty"`
}

// Basic represents the pushgateway basic auth config
type Basic struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

// Bearer represents the pushgateway bearer auth config
type Bearer struct {
	Token string `json:"token"`
}
