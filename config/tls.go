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
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type (
	// TLS holds the configuration settings for a TLS Cert.
	TLS struct {
		cert   string
		key    string
		cacert string
	}
)

const (
	tlsPrefix = "tls"
)

// NewTLS returns the tls configuration.
func NewTLS(opts ...Option) (TLS, error) {
	o := &Options{}
	for _, opt := range opts {
		opt.Apply(o)
	}
	prefix := o.Prefix + tlsPrefix

	config := struct {
		Cert   string
		Key    string
		CACert string
	}{}
	if err := envconfig.Process(prefix, &config); err != nil {
		return TLS{}, fmt.Errorf("failed to load %s configuration: %w", prefix, err)
	}

	return TLS{
		cert:   strings.TrimSpace(config.Cert),
		key:    strings.TrimSpace(config.Key),
		cacert: strings.TrimSpace(config.CACert),
	}, nil
}

// Cert returns the path of the certificate file. The value is set from the
// <PREFIX_>TLS_CERT environment variable.
func (t TLS) Cert() string {
	return t.cert
}

// Key returns the path of the certificate key file. The value is set from the
// <PREFIX_>TLS_KEY environment variable.
func (t TLS) Key() string {
	return t.key
}

// CACert returns the path of the certificate of the CA certificate file. This
// is used when creating a TLS connection with an entity that is presenting a
// certificate that is not signed by a well known CA available in the OS CA
// bundle. The value is set from the <PREFIX_>TLS_CACERT environment
// variable.
func (t TLS) CACert() string {
	return t.cacert
}

// TLSConfig will create a *tls.Config given the options. This will return an error
// if there is a problem loading the required certificate files.  If the
// WithMTLS option is specified, a client CA cert is required.
func (t TLS) TLSConfig(opts ...TLSOption) (*tls.Config, error) {
	cfg := &tls.Config{}
	for _, opt := range opts {
		opt.Apply(cfg)
	}

	// Load the tls certificate.
	cert, err := tls.LoadX509KeyPair(t.cert, t.key)
	if err != nil {
		return nil, fmt.Errorf("failed to load tls certificate: %w", err)
	}
	cfg.Certificates = append(cfg.Certificates, cert)

	// If we are doing mTLS...
	if cfg.ClientAuth == tls.RequireAndVerifyClientCert {
		// .. and we have a CA cert
		caCertCfg := t.cacert
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
	// TLSOption provides options for configuring the creation of a tls.Config.
	TLSOption interface {
		Apply(*tls.Config)
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

func (o tlsOption) Apply(cfg *tls.Config) {
	o.f(cfg)
}
