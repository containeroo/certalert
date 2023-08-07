package config

import "certalert/internal/certificates"

// DeepCopy returns a deep copy of the config
func (c Config) DeepCopy() Config {
	// Copying basic types (like string, int, bool)
	// is straightforward since they don't contain internal references.
	var newConfig Config
	newConfig.Server = Server{
		Hostname: c.Server.Hostname,
		Port:     c.Server.Port,
	}
	newConfig.Pushgateway = Pushgateway{
		Address:            c.Pushgateway.Address,
		InsecureSkipVerify: c.Pushgateway.InsecureSkipVerify,
		Job:                c.Pushgateway.Job,
		Auth: Auth{
			Basic: Basic{
				Username: c.Pushgateway.Auth.Basic.Username,
				Password: c.Pushgateway.Auth.Basic.Password,
			},
			Bearer: Bearer{
				Token: c.Pushgateway.Auth.Bearer.Token,
			},
		},
	}

	// For slices, you'll want to ensure you're creating a new slice
	// and copying each element (especially if they are structs).
	newCerts := make([]certificates.Certificate, len(c.Certs))
	for i, cert := range c.Certs {
		newCerts[i] = certificates.Certificate{
			Name:     cert.Name,
			Enabled:  cert.Enabled, // Assuming bool pointers are okay to copy directly
			Path:     cert.Path,
			Password: cert.Password,
			Type:     cert.Type,
		}
	}
	newConfig.Certs = newCerts

	return newConfig
}
