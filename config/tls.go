// Copyright 2021-2022 arcadium.dev <info@arcadium.dev>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package config // import "arcadium.dev/core/config"

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
)

// NewTLS will create a *tls.Config given the config and options. This will
// return an error if there is a problem loading the required certificate files.
// If the WithMTLS option is specified, a client CA cert is required.
func NewTLS(config TLS, opts ...TLSOption) (*tls.Config, error) {
	cfg := &tls.Config{}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	// Load the server certificate.
	cert, err := tls.LoadX509KeyPair(config.Cert(), config.Key())
	if err != nil {
		return nil, fmt.Errorf("failed to load server certificate: %w", err)
	}
	cfg.Certificates = append(cfg.Certificates, cert)

	// If we are doing mTLS...
	if cfg.ClientAuth == tls.RequireAndVerifyClientCert {
		// .. and we have a CA cert
		caCertCfg := config.CACert()
		if caCertCfg != "" {
			// ... create a new, empty CA certificate pool and add client's CA cert to it.
			cfg.ClientCAs = x509.NewCertPool()
			caCert, err := os.ReadFile(caCertCfg)
			if err != nil {
				return nil, fmt.Errorf("failed to load the client CA certificate: %w", err)
			}
			cfg.ClientCAs.AppendCertsFromPEM(caCert)
		}
	}

	return cfg, nil
}

type (
	// TLS contains the information necessary to create a tls.Config.
	TLS interface {
		// Cert returns the file name of the PEM encoded public key.
		Cert() string

		// Key returns the file name of the PEM encoded private key.
		Key() string

		// CACert returns the file name of the PEM encoded public key of the client CA.
		CACert() string
	}

	// TLSOption provides options for configuring the creation of a tls.Config.
	TLSOption interface {
		apply(*tls.Config)
	}
)

// WithMTLS will setup the tls.Config to require and verify client connections.
func WithMTLS() TLSOption {
	return newTLSOption(func(cfg *tls.Config) {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	})
}

type (
	tlsOption struct {
		f func(*tls.Config)
	}
)

func newTLSOption(f func(*tls.Config)) tlsOption {
	return tlsOption{f: f}
}

func (o tlsOption) apply(cfg *tls.Config) {
	o.f(cfg)
}
