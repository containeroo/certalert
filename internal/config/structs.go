package config

import "certalert/internal/certificates"

// Config represents the config file
var App Config

// Config represents the config file
type Config struct {
	Server      Server                     `yaml:"server"`
	Pushgateway Pushgateway                `yaml:"pushgateway"`
	Certs       []certificates.Certificate `yaml:"certs"`
}

// Server represents the server config
type Server struct {
	Hostname string `yaml:"hostname"`
	Port     int    `yaml:"port"`
}

// Pushgateway represents the pushgateway config
type Pushgateway struct {
	Address    string `yaml:"address"`
	SkipVerify bool   `yaml:"skip_verify"`
	Job        string `yaml:"job"`
	Auth       Auth   `yaml:"auth"`
}

// Auth represents the pushgateway auth config
type Auth struct {
	Basic  Basic  `yaml:"basic,omitempty"`
	Bearer Bearer `yaml:"bearer,omitempty"`
}

// Basic represents the pushgateway basic auth config
type Basic struct {
	Password string `yaml:"password"`
	Username string `yaml:"username"`
}

// Bearer represents the pushgateway bearer auth config
type Bearer struct {
	Token string `yaml:"token"`
}
