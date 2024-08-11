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
	"fmt"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
)

const (
	defaultAddr            = ":8443"
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 10 * time.Second
	defaultShutdownTimeout = 10 * time.Second
)

type (
	// Server represents an HTTP server.
	Server struct {
		addr            string
		corsOptions     *cors.Options
		shutdownTimeout time.Duration

		logger   *zerolog.Logger
		listener net.Listener
		server   *http.Server
		router   *mux.Router
		scheme   string

		mu       sync.RWMutex
		services []Service
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
		Shutdown()
	}

	corsLogger struct {
		logger *zerolog.Logger
	}
)

func (c corsLogger) Printf(f string, v ...any) {
	c.logger.Debug().Msgf(f, v...)
}

// New creates an HTTP server with and has not started to accept requests yet.
func New(ctx context.Context, opts ...Option) *Server {
	s := &Server{
		addr:   defaultAddr,
		logger: zerolog.Ctx(ctx),
		server: &http.Server{
			ReadTimeout:  defaultReadTimeout,
			WriteTimeout: defaultWriteTimeout,
		},
		router:          mux.NewRouter(),
		scheme:          "http",
		shutdownTimeout: defaultShutdownTimeout,
	}

	// Load options.
	for _, opt := range opts {
		opt.apply(s)
	}

	// The CORS handler needs to be invoked before the mux so that the CORS handler
	// can handle preflight requests before they would hit the mux. Otherwise the
	// mux would try to route those requests.
	var c *cors.Cors
	if s.corsOptions == nil {
		s.logger.Info().Msg("cors allow all")
		c = cors.AllowAll()
	} else {
		s.logger.Info().Msgf("cors allowed origins: %q", s.corsOptions.AllowedOrigins)
		s.logger.Info().Msgf("cors allowed methods: %q", s.corsOptions.AllowedMethods)
		s.logger.Info().Msgf("cors allowed headers: %q", s.corsOptions.AllowedHeaders)
		c = cors.New(*s.corsOptions)
	}
	c.Log = corsLogger{logger: s.logger}

	s.server.Handler = c.Handler(s.router)

	// Set up the logging fields.
	tlsMsg := ""
	if s.server.TLSConfig != nil {
		s.scheme = "https"
		tlsMsg = ", tls: enabled"
		if s.server.TLSConfig.ClientAuth == tls.RequireAndVerifyClientCert {
			tlsMsg = ", mtls: enabled"
		}
	}
	s.logger.Info().Msgf("%s server created, address '%s'%s", s.scheme, s.addr, tlsMsg)

	return s
}

// Middleware installs the given middleware with the router.
func (s *Server) Middleware(mw ...mux.MiddlewareFunc) {
	if len(mw) > 0 {
		s.router.Use(mw...)
	}
}

// Register associates the given services with the router.
func (s *Server) Register(services ...Service) {
	s.mu.Lock()
	s.services = append(s.services, services...)
	s.mu.Unlock()

	r := s.router.PathPrefix("/").Subrouter()
	for _, service := range services {
		service.Register(r)
		s.logger.Info().Msgf("http service registered: %s", service.Name())
	}
}

// Serve accepts incoming connections, creating a new service goroutine for each. The
// service goroutine reads requests and then call the handler to reply to them.
func (s *Server) Serve() error {
	var err error
	if s.listener, err = net.Listen("tcp", s.addr); err != nil {
		return fmt.Errorf("failed to listen on address '%s', %w", s.addr, err)
	}

	serviceNames := make([]string, 0)
	s.mu.RLock()
	for _, service := range s.services {
		serviceNames = append(serviceNames, service.Name())
	}
	s.mu.RUnlock()
	services := strings.Join(serviceNames, ",")

	s.logger.Info().Msgf("begin serving %s, address '%s', services: %s", s.scheme, s.addr, services)
	defer s.logger.Info().Msgf("serving %s complete, address '%s', services: %s", s.scheme, s.addr, services)

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
func (s *Server) Shutdown() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(s.shutdownTimeout))
	defer cancel()

	// Stop each service.
	s.mu.RLock()
	for _, service := range s.services {
		service.Shutdown()
		s.logger.Info().Msgf("http service shutdown, service: %s", service.Name())
	}
	s.mu.RUnlock()

	// Stop the http server.
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Err(err).Msg("failed to shutdown http server")
	}

	s.logger.Info().Msg("http server shutdown")
}
