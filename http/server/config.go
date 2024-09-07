// Copyright 2021-2024 arcadium.dev <info@arcadium.dev>
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

package server // import "arcadium.dev/http/server"

import (
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type (
	// Config holds the configuration information for the restful api server.
	Config struct {
		serverAddr      string
		tlsCert         string
		tlsKey          string
		tlsCACert       string
		mtlsEnabled     bool
		allowedOrigins  []string
		allowedMethods  []string
		allowedHeaders  []string
		readTimeout     time.Duration
		writeTimeout    time.Duration
		shutdownTimeout time.Duration
		pprofEnabled    bool
	}
)

// NewConfig returns the configuration of restful api server.
func NewConfig(prefix ...string) (Config, error) {
	cfg := struct {
		ServerAddr      string        `required:"true" split_words:"true"`
		TlsCert         string        `split_words:"true"`
		TlsKey          string        `split_words:"true"`
		TlsCacert       string        `split_words:"true"`
		MtlsEnabled     bool          `split_words:"true"`
		AllowedOrigins  string        `split_words:"true"`
		AllowedMethods  string        `split_words:"true"`
		AllowedHeaders  string        `split_words:"true"`
		ReadTimeout     time.Duration `split_words:"true"`
		WriteTimeout    time.Duration `split_words:"true"`
		ShutdownTimeout time.Duration `split_words:"true"`
		PprofEnabled    bool          `split_words:"true"`
	}{}

	pfix := ""
	if len(prefix) == 1 {
		pfix = prefix[0]
	}
	if err := envconfig.Process(pfix, &cfg); err != nil {
		return Config{}, fmt.Errorf("failed to load configuration: %w", err)
	}

	c := Config{
		serverAddr:      strings.TrimSpace(cfg.ServerAddr),
		tlsCert:         strings.TrimSpace(cfg.TlsCert),
		tlsKey:          strings.TrimSpace(cfg.TlsKey),
		tlsCACert:       strings.TrimSpace(cfg.TlsCacert),
		mtlsEnabled:     cfg.MtlsEnabled,
		readTimeout:     cfg.ReadTimeout,
		writeTimeout:    cfg.WriteTimeout,
		shutdownTimeout: cfg.ShutdownTimeout,
		pprofEnabled:    cfg.PprofEnabled,
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

// ServerAddr returns the network address the server will listen on.
func (c Config) ServerAddr() string {
	return c.serverAddr
}

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

// ReadTimeout returns the server read timeout.
func (c Config) ReadTimeout() time.Duration {
	return c.readTimeout
}

// WriteTimeout returns the server write timeout.
func (c Config) WriteTimeout() time.Duration {
	return c.writeTimeout
}

// ShutdownTimeout returns the server shutdown timeout.
func (c Config) ShutdownTimeout() time.Duration {
	return c.shutdownTimeout
}

// PprofEnabled returns true when pprof should be enabled.
func (c Config) PprofEnabled() bool {
	return c.pprofEnabled
}

// ToOptions converts the configuration to a list of options.
func (c Config) ToOptions() []Option {
	var options []Option
	if c.serverAddr != "" {
		options = append(options, WithAddr(c.serverAddr))
	}
	if c.tlsCert != "" && c.tlsKey != "" {
		options = append(options, WithTLSCert(c.tlsCert, c.tlsKey))
	}
	if c.tlsCACert != "" {
		options = append(options, WithTLSClientCACert(c.tlsCACert))
	}
	if c.mtlsEnabled {
		options = append(options, WithMTLSEnabled(c.mtlsEnabled))
	}
	if len(c.allowedOrigins) > 0 {
		options = append(options, WithCORSAllowedOrigins(c.allowedOrigins))
	}
	if len(c.allowedMethods) > 0 {
		options = append(options, WithCORSAllowedMethods(c.allowedMethods))
	}
	if len(c.allowedHeaders) > 0 {
		options = append(options, WithCORSAllowedHeaders(c.allowedHeaders))
	}
	if c.writeTimeout > 0 {
		options = append(options, WithWriteTimeout(c.writeTimeout))
	}
	if c.readTimeout > 0 {
		options = append(options, WithReadTimeout(c.readTimeout))
	}
	if c.shutdownTimeout > 0 {
		options = append(options, WithShutdownTimeout(c.shutdownTimeout))
	}
	return options
}
