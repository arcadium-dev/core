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

package http // import "arcadium.dev/core/http"

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"runtime"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/mux"

	"arcadium.dev/core/log"
)

const (
	defaultAddr            = ":8443"
	defaultShutdownTimeout = 10 * time.Second
)

var (
	defaultLogger log.Logger
)

func init() {
	defaultLogger, _ = log.New(log.WithFormat(log.FormatNop))
}

type (
	// Server represents an HTTP server.
	Server struct {
		addr            string
		shutdownTimeout time.Duration

		logger     log.Logger
		listener   net.Listener
		server     *http.Server
		router     *mux.Router
		middleware []mux.MiddlewareFunc

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
)

// NewServer creates an HTTP server with and has not started to accept requests yet.
func NewServer(opts ...ServerOption) *Server {
	s := &Server{
		addr:            defaultAddr,
		logger:          defaultLogger,
		server:          &http.Server{},
		router:          mux.NewRouter(),
		shutdownTimeout: defaultShutdownTimeout,
	}
	s.server.Handler = s.router

	// Load options.
	for _, opt := range opts {
		opt.apply(s)
	}

	// Set up the logging fields.
	msg := []interface{}{
		"msg", "http server created",
		"addr", s.addr,
	}
	if s.server.TLSConfig != nil {
		if s.server.TLSConfig.ClientAuth == tls.RequireAndVerifyClientCert {
			msg = append(msg, "mtls", "enabled")
		} else {
			msg = append(msg, "tls", "enabled")
		}
	}
	s.logger.Info(msg...)

	s.router.Use(s.recoverPanics)
	s.router.Use(s.requestLogging)
	if len(s.middleware) > 0 {
		s.router.Use(s.middleware...)
	}

	return s
}

// Register associates the given services with the router.
func (s *Server) Register(services ...Service) {
	s.mu.Lock()
	s.services = append(s.services, services...)
	s.mu.Unlock()

	r := s.router.PathPrefix("/").Subrouter()
	for _, service := range services {
		service.Register(r)
		s.logger.Info("msg", "service registered", "service", service.Name())
	}
}

// Serve accepts incoming connections, creating a new service goroutine for each. The
// service goroutine reads requests and then call the handler to reply to them.
func (s *Server) Serve() error {
	var err error
	if s.listener, err = net.Listen("tcp", s.addr); err != nil {
		return fmt.Errorf("failed to listen on %s: %w", s.addr, err)
	}

	serviceNames := make([]string, 0)
	s.mu.RLock()
	for _, service := range s.services {
		serviceNames = append(serviceNames, service.Name())
	}
	s.mu.RUnlock()
	services := strings.Join(serviceNames, ",")

	s.logger.Info("msg", "begin serving", "services", services, "addr", s.addr)
	defer s.logger.Info("msg", "serving complete", "services", services, "addr", s.addr)

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
		s.logger.Info("msg", "service shutdown", "service", service.Name())
	}
	s.mu.RUnlock()

	// Stop the http server.
	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("msg", "failed to shutdown", "error", err.Error())
	}

	s.logger.Info("msg", "infra shutdown")
}

// recoverPanics is middleware for recovering and reporting panics.
func (s *Server) recoverPanics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Install the recovery and reporting function.
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				msg := []interface{}{
					"msg", "recovering from a panic",
					"req", fmt.Sprintf("%+v", *r),
				}
				s.logger.Error(msg...)

				buf := make([]byte, 4096)
				n := runtime.Stack(buf, false)
				buf = buf[:n]
				s.logger.Error("stacktrace", fmt.Sprintf("%s", string(buf)))
			}
		}()
		// Delegate to next handler in middleware chain.
		next.ServeHTTP(w, r)
	})
}

// requestLogging is middleware to create a request specific logger, passing
// in through the request's context, as well as log the incoming request.
func (s *Server) requestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fields := []interface{}{
			"method", r.Method,
			"url", r.URL.String(),
		}
		l := s.logger.With(fields...)
		req := r.Clone(log.NewContextWithLogger(r.Context(), l))

		l.Debug("msg", "request start")
		next.ServeHTTP(w, req)
		l.Debug("msg", "request complete")
	})
}
