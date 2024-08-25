// Copyright 2024 arcadium.dev <info@arcadium.dev>
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

package ssh // import "arcadium.dev/core/ssh"

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/rs/cors"

	"arcadium.dev/core/build"
	"arcadium.dev/core/http/server"
	"arcadium.dev/core/http/services"
)

type (
	HTTPServer struct {
		server   *server.Server
		services []server.Service

		cfg    Config
		wg     *sync.WaitGroup
		result chan error
	}
)

// NewHTTPServer returns a new http server.
func NewHTTPServer(ctx context.Context, cfg Config) (*HTTPServer, error) {
	var tlsConfig *tls.Config

	// Setup HTTPS.
	if cfg.TLSCert() != "" && cfg.TLSKey() != "" {
		cert, err := tls.LoadX509KeyPair(cfg.TLSCert(), cfg.TLSKey())
		if err != nil {
			return nil, fmt.Errorf("failed to load tls certificate: %w", err)
		}
		tlsConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}
	}
	if cfg.TLSCACert() != "" {
		tlsConfig.ClientCAs = x509.NewCertPool()
		caCert, err := os.ReadFile(cfg.TLSCACert())
		if err != nil {
			return nil, fmt.Errorf("failed to load the tls client CA certificate: %w", err)
		}
		tlsConfig.ClientCAs.AppendCertsFromPEM(caCert)
	}

	// Gather the server options.
	var opts []server.Option
	opts = append(opts,
		server.WithAddr(cfg.HTTPServerAddr()),
		server.WithTLS(tlsConfig),
	)

	// Setup CORS.
	corsOpts := &cors.Options{}
	if len(cfg.AllowedOrigins()) != 0 {
		corsOpts.AllowedOrigins = cfg.AllowedOrigins()
	}
	if len(cfg.AllowedMethods()) != 0 {
		corsOpts.AllowedMethods = cfg.AllowedMethods()
	}
	if len(cfg.AllowedHeaders()) != 0 {
		corsOpts.AllowedHeaders = cfg.AllowedHeaders()
	}
	if len(corsOpts.AllowedOrigins) == 1 && corsOpts.AllowedOrigins[0] != "*" {
		corsOpts.AllowCredentials = true
	}
	if len(corsOpts.AllowedOrigins) > 0 || len(corsOpts.AllowedMethods) > 0 || len(corsOpts.AllowedHeaders) > 0 {
		opts = append(opts, server.WithCORS(corsOpts))
	}

	return &HTTPServer{
		server: server.New(ctx, opts...),
		cfg:    cfg,
		wg:     &sync.WaitGroup{},
		result: make(chan error, 1),
	}, nil
}

func (s HTTPServer) Start(info build.Information) {
	s.wg.Add(1)

	s.services = []server.Service{
		services.Health{Start: time.Now(), Info: info},
		services.Metrics{},
	}
	if s.cfg.PProfEnabled() {
		s.services = append(s.services, services.PProf{})
	}
	s.server.Register(s.services...)

	// Serve.
	go func() {
		s.wg.Done()
		s.result <- s.server.Serve()
	}()
}

func (s *HTTPServer) Done() <-chan error {
	return s.result
}

func (s HTTPServer) Shutdown() {
	s.wg.Wait()
	s.server.Shutdown()
}
