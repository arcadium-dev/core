// Copyright 2021 arcadium.dev <info@arcadium.dev>
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

package server // import "arcadium.dev/core/server"

//go:generate mockgen -package mockserver -destination ./mock/config.go . Config

import (
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/pkg/errors"
)

// Config contains the information necessary to create a server.
type Config interface {

	// Addr returns that network address the server listens to.
	Addr() string

	// Cert returns the file name of the PEM encoded public key.
	Cert() string

	// Key returns the file name of the PEM encoded private key.
	Key() string

	// CACert returns the file name of the PEM encoded public key of the client CA.
	CACert() string
}

// CreateTLSConfig will create a tls.Config given the config and options. This will
// return an error if there is a problem loading the required certificate files.
// If the WithMTLS option is specified, a client CA cert is required.
func CreateTLSConfig(config Config, opts ...TLSConfigOption) (*tls.Config, error) {
	cfg := &tls.Config{}
	for _, opt := range opts {
		opt.apply(cfg)
	}

	// Load the server certificate.
	cert, err := tls.LoadX509KeyPair(config.Cert(), config.Key())
	if err != nil {
		return nil, errors.Wrap(err, "Failed to load server certificate")
	}
	cfg.Certificates = append(cfg.Certificates, cert)

	// If we are doing mTLS...
	if cfg.ClientAuth == tls.RequireAndVerifyClientCert {
		file := config.CACert()
		if file == "" {
			return nil, errors.New("A client CA certification must be configured for mTLS.")
		}
		// ... create a CA certificate pool and add client's CA cert to it.
		cfg.ClientCAs = x509.NewCertPool()
		caCert, err := os.ReadFile(file)
		if err != nil {
			return nil, errors.Wrap(err, "Failed to load the client CA certificate")
		}
		cfg.ClientCAs.AppendCertsFromPEM(caCert)
	}

	return cfg, nil
}

// WithMTLS will setup the tls.Config to require and verify client connections.
func WithMTLS() TLSConfigOption {
	return newTLSConfigOption(func(cfg *tls.Config) {
		cfg.ClientAuth = tls.RequireAndVerifyClientCert
	})
}

type (
	// TLSConfigOption provides options for configuring the creation of a tls.Config.
	TLSConfigOption interface {
		apply(*tls.Config)
	}

	tlsConfigOption struct {
		f func(*tls.Config)
	}
)

func newTLSConfigOption(f func(*tls.Config)) tlsConfigOption {
	return tlsConfigOption{f: f}
}

func (o tlsConfigOption) apply(cfg *tls.Config) {
	o.f(cfg)
}
