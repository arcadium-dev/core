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

// Package server provides an http server.
package server // import "arcadium.dev/core/http/server"

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

const (
	defaultAddr            = ":80"
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 10 * time.Second
	defaultShutdownTimeout = 10 * time.Second
)

type (
	// Server represents an HTTP server.
	Server struct {
		addr            string
		shutdownTimeout time.Duration

		tlsConfig   *tls.Config
		tlsCert     string
		tlsKey      string
		tlsCACert   string
		mtlsEnabled bool

		corsOptions        *cors.Options
		corsAllowedOrigins []string
		corsAllowedMethods []string
		corsAllowedHeaders []string

		listener net.Listener
		server   *http.Server
		router   *mux.Router
		scheme   string

		services *services
	}

	// Service defines the methods required by the Server to register with
	// the service with the router.
	Service interface {
		// Register will register this service with the given router.
		Register(router *mux.Router)

		// Name provides the name of the service.
		Name() string

		// Shutdown allows the service to stop any long running background processes it
		// may have.
		Shutdown(context.Context)
	}

	services struct {
		mu       sync.RWMutex
		services []Service
	}

	corsLogger struct {
		logger *zerolog.Logger
	}
)

func (c corsLogger) Printf(f string, v ...any) {
	c.logger.Debug().Msgf(f, v...)
}

func (s *services) append(srv []Service) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.services = append(s.services, srv...)
}

func (s *services) len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.services)
}

func (s *services) index(i int) Service {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.services[i]
}

// New creates an HTTP server with and has not started to accept requests yet.
func New(ctx context.Context, opts ...Option) (*Server, error) {
	logger := zerolog.Ctx(ctx)

	s := &Server{
		addr: defaultAddr,
		server: &http.Server{
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
		},
		router:          mux.NewRouter(),
		scheme:          "http",
		shutdownTimeout: defaultShutdownTimeout,
		services:        &services{},
	}

	// Load options.
	for _, opt := range opts {
		opt.apply(s)
	}

	if err := s.setupTLS(ctx); err != nil {
		return nil, err
	}
	s.setupCORS(ctx)

	// Set up the logging fields.
	tlsMsg := ""
	if s.server.TLSConfig != nil {
		s.scheme = "https"
		tlsMsg = ", tls: enabled"
		if s.server.TLSConfig.ClientAuth == tls.RequireAndVerifyClientCert {
			tlsMsg = ", mtls: enabled"
		}
	}
	logger.Info().Msgf("%s server created, address '%s'%s", s.scheme, s.addr, tlsMsg)

	return s, nil
}

func (s *Server) setupTLS(ctx context.Context) error {
	// If the entire tls.Config was given, prefer that.
	if s.tlsConfig != nil {
		if s.tlsCert != "" || s.tlsKey != "" || s.tlsCACert != "" {
			zerolog.Ctx(ctx).Warn().Msg("both the TLS config and individual tls properties were defined, using TLS config")
		}
		s.server.TLSConfig = s.tlsConfig
		return nil
	}

	// Ensure both tlsCert and tlsKey are defined when one is defined.
	switch {
	case s.tlsCert != "" && s.tlsKey == "":
		return fmt.Errorf("the tls key must be defined with then tls cert is defined")
	case s.tlsCert == "" && s.tlsKey != "":
		return fmt.Errorf("the tls cert must be defined with then tls key is defined")
	}

	var tlsConfig *tls.Config
	if s.tlsCert != "" && s.tlsKey != "" {
		// Setup cert used for HTTPS.
		cert, err := tls.LoadX509KeyPair(s.tlsCert, s.tlsKey)
		if err != nil {
			return fmt.Errorf("failed to load TLS certificate: %w", err)
		}
		tlsConfig = &tls.Config{Certificates: []tls.Certificate{cert}}

		// Setup the client CA cert if available.
		if s.tlsCACert != "" {
			tlsConfig.ClientCAs = x509.NewCertPool()
			clientCACert, err := os.ReadFile(s.tlsCACert)
			if err != nil {
				return fmt.Errorf("failed to load the TLS client CA certificate: %w", err)
			}
			tlsConfig.ClientCAs.AppendCertsFromPEM(clientCACert)
		}

		// Setup MTLS if enabled.
		if s.mtlsEnabled {
			tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert
		}
	}
	s.server.TLSConfig = tlsConfig

	return nil
}

func (s *Server) setupCORS(ctx context.Context) {
	if s.corsOptions != nil {
		if len(s.corsAllowedOrigins) > 0 || len(s.corsAllowedMethods) > 0 || len(s.corsAllowedHeaders) > 0 {
			zerolog.Ctx(ctx).Warn().Msg("both the cors options and individual cors properties were defined, using cors options")
		}
		s.finishCORS(ctx, s.corsOptions)
		return
	}

	if len(s.corsAllowedOrigins) == 0 && len(s.corsAllowedMethods) == 0 && len(s.corsAllowedHeaders) == 0 {
		s.finishCORS(ctx, nil)
		return
	}

	corsOpts := &cors.Options{}
	if len(s.corsAllowedOrigins) > 0 {
		corsOpts.AllowedOrigins = s.corsAllowedOrigins
	}
	if len(s.corsAllowedMethods) > 0 {
		corsOpts.AllowedMethods = s.corsAllowedMethods
	}
	if len(s.corsAllowedHeaders) > 0 {
		corsOpts.AllowedHeaders = s.corsAllowedHeaders
	}
	if len(corsOpts.AllowedOrigins) == 1 && corsOpts.AllowedOrigins[0] != "*" {
		corsOpts.AllowCredentials = true
	}
	s.finishCORS(ctx, corsOpts)
}

func (s *Server) finishCORS(ctx context.Context, opts *cors.Options) {
	logger := zerolog.Ctx(ctx)

	// The CORS handler needs to be invoked before the mux so that the CORS handler
	// can handle preflight requests before they would hit the mux. Otherwise the
	// mux would try to route those requests.
	var c *cors.Cors
	if opts == nil {
		logger.Info().Msg("cors allow all")
		c = cors.AllowAll()
	} else {
		logger.Info().Msgf("cors allowed origins: %q", opts.AllowedOrigins)
		logger.Info().Msgf("cors allowed methods: %q", opts.AllowedMethods)
		logger.Info().Msgf("cors allowed headers: %q", opts.AllowedHeaders)
		c = cors.New(*opts)
	}
	c.Log = corsLogger{logger: logger}

	s.server.Handler = c.Handler(s.router)
}

// Middleware installs the given middleware with the router.
func (s Server) Middleware(mw ...mux.MiddlewareFunc) {
	if len(mw) > 0 {
		s.router.Use(mw...)
	}
}

// Register associates the given services with the router.
func (s *Server) Register(ctx context.Context, services ...Service) {
	s.services.append(services)

	r := s.router.PathPrefix("/").Subrouter()
	for _, service := range services {
		service.Register(r)
		zerolog.Ctx(ctx).Info().Msgf("http service registered: %s", service.Name())
	}
}

// Serve accepts incoming connections. This is a blocking call and should be
// called in the context of a new goroutime.
func (s Server) Serve(ctx context.Context) error {
	var err error
	if s.listener, err = net.Listen("tcp", s.addr); err != nil {
		return fmt.Errorf("failed to listen on address '%s', %w", s.addr, err)
	}

	serviceNames := make([]string, 0)
	for i := 0; i < s.services.len(); i++ {
		service := s.services.index(i)
		serviceNames = append(serviceNames, service.Name())
	}
	services := strings.Join(serviceNames, ",")

	zerolog.Ctx(ctx).Info().Msgf("begin serving %s, address '%s', services: %s", s.scheme, s.addr, services)
	defer zerolog.Ctx(ctx).Info().Msgf("serving %s complete, address '%s', services: %s", s.scheme, s.addr, services)

	if s.server.TLSConfig != nil {
		err = s.server.ServeTLS(s.listener, "", "")
	} else {
		err = s.server.Serve(s.listener)
	}

	if err == http.ErrServerClosed {
		err = nil
	}
	return err
}

// Shutdown stops the http server gracefully without interrupting any active connections.
// It will, however, forcefully stop if the shutdown timeout expires while shutting down.
func (s Server) Shutdown(ctx context.Context) {
	ctx, cancel := context.WithDeadline(ctx, time.Now().Add(s.shutdownTimeout))
	defer cancel()

	// Stop each service.
	for i := 0; i < s.services.len(); i++ {
		service := s.services.index(i)
		service.Shutdown(ctx)
		zerolog.Ctx(ctx).Info().Msgf("http service shutdown, service: %s", service.Name())
	}

	// Stop the http server.
	if err := s.server.Shutdown(ctx); err != nil {
		zerolog.Ctx(ctx).Err(err).Msg("failed to shutdown http server")
	}

	zerolog.Ctx(ctx).Info().Msg("http server shutdown")
}
