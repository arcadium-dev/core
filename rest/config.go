// Copyright 2021-2023 arcadium.dev <info@arcadium.dev>
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

package rest // import "arcadium.dev/core/rest"

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type (
	// Config holds the configuration information for the restful api server.
	Config struct {
		logLevel       string
		tlsCert        string
		tlsKey         string
		tlsCACert      string
		mtlsEnabled    bool
		serverAddr     string
		allowedOrigins []string
		allowedMethods []string
		allowedHeaders []string
		pprofEnabled   bool
	}
)

const (
	DefaultLogLevel = "info"
)

// NewConfig returns the configuration of restful api server.
func NewConfig() (Config, error) {
	cfg := struct {
		LogLevel       string `split_words:"true"`
		TlsCert        string `split_words:"true"`
		TlsKey         string `split_words:"true"`
		TlsCacert      string `split_words:"true"`
		MtlsEnabled    bool   `split_words:"true"`
		ServerAddr     string `required:"true" split_words:"true"`
		AllowedOrigins string `split_words:"true"`
		AllowedMethods string `split_words:"true"`
		AllowedHeaders string `split_words:"true"`
		PprofEnabled   bool   `split_words:"true"`
	}{
		LogLevel: DefaultLogLevel,
	}
	if err := envconfig.Process("", &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	c := Config{
		logLevel:     strings.TrimSpace(strings.ToLower(cfg.LogLevel)),
		tlsCert:      strings.TrimSpace(cfg.TlsCert),
		tlsKey:       strings.TrimSpace(cfg.TlsKey),
		tlsCACert:    strings.TrimSpace(cfg.TlsCacert),
		mtlsEnabled:  cfg.MtlsEnabled,
		serverAddr:   strings.TrimSpace(cfg.ServerAddr),
		pprofEnabled: cfg.PprofEnabled,
	}

	origins := strings.TrimSpace(cfg.AllowedOrigins)
	if len(origins) > 0 {
		for _, o := range strings.Split(origins, ",") {
			c.allowedOrigins = append(c.allowedOrigins, strings.TrimSpace(o))
		}
	}
	methods := strings.TrimSpace(cfg.AllowedMethods)
	if len(methods) > 0 {
		for _, m := range strings.Split(methods, ",") {
			c.allowedMethods = append(c.allowedMethods, strings.TrimSpace(m))
		}
	}
	headers := strings.TrimSpace(cfg.AllowedHeaders)
	if len(headers) > 0 {
		for _, h := range strings.Split(headers, ",") {
			c.allowedHeaders = append(c.allowedHeaders, strings.TrimSpace(h))
		}
	}

	return c, nil
}

// LogLevel returns the logging level. The default level is "error".
func (c Config) LogLevel() string { return c.logLevel }

// TLSCert returns the path of the certificate file.
func (c Config) TLSCert() string {
	return c.tlsCert
}

// TLSKey returns the path of the certificate key file.
func (c Config) TLSKey() string {
	return c.tlsKey
}

// TLSCACert returns the path of the CA certificate file.
func (c Config) TLSCACert() string {
	return c.tlsCACert
}

// MTLSEnabled returns true if MTLS is enabled.
func (c Config) MTLSEnabled() bool {
	return c.mtlsEnabled
}

// ServerAddr returns the network address the server will listen on.
func (c Config) ServerAddr() string {
	return c.serverAddr
}

// AllowedOrigins returns a string of Allowed Origins passed to the CORS middleware.
func (c Config) AllowedOrigins() []string {
	return c.allowedOrigins
}

// AllowedMethods returns a string of Allowed Methods passed to the CORS middleware.
func (c Config) AllowedMethods() []string {
	return c.allowedMethods
}

// AllowedHeaders returns a string of Allowed Headers passed to the CORS middleware.
func (c Config) AllowedHeaders() []string {
	return c.allowedHeaders
}

// PProfEnabled adds the pprof endpoints to the server.
func (c Config) PProfEnabled() bool {
	return c.pprofEnabled
}
