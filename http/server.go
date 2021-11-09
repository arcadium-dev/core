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

package http // import "arcadium.dev/core/server/http"

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	_ "net/http/pprof" // Provide pprof profiling data.
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"arcadium.dev/core/config"
	"arcadium.dev/core/errors"
	"arcadium.dev/core/log"
)

const (
	defaultShutdownTimeout = 10 * time.Second
)

type (
	// Server represents an HTTP server.
	Server struct {
		addr            string
		shutdownTimeout time.Duration

		requestCount   *prometheus.CounterVec
		requestSeconds *prometheus.CounterVec

		logger   log.Logger
		listener net.Listener
		server   *http.Server
		router   *mux.Router
	}

	// Service defines the methods required by the Server to register with
	// the service with the router.
	Service interface {
		// Register will register this service with the given router.
		Register(router *mux.Router)

		// Fields provides a set of fields for logging.
		Fields() []interface{}
	}

	// Config contains the information necessary to create an HTTP server
	Config interface {
		config.TLS

		// Addr returns that network address the server listens to.
		Addr() string
	}
)

// New creates an HTTP server with and has not started to accept requests yet.
func New(cfg Config, logger log.Logger, opts ...Option) *Server {
	s := &Server{
		addr:            cfg.Addr(),
		logger:          logger,
		server:          &http.Server{},
		router:          mux.NewRouter(),
		shutdownTimeout: defaultShutdownTimeout,
	}

	// Load options.
	for _, opt := range opts {
		opt.apply(s)
	}

	// Set up the logging fields.
	fields := []interface{}{
		"addr", s.addr,
	}
	if s.server.TLSConfig != nil {
		if s.server.TLSConfig.ClientAuth == tls.RequireAndVerifyClientCert {
			fields = append(fields, "mtls", "enabled")
		} else {
			fields = append(fields, "tls", "enabled")
		}
	}
	s.logger = s.logger.With(fields...)

	// Setup middleware.
	s.router.Use(s.requestLogging)
	s.router.Use(s.requestMetrics)
	s.router.Use(s.catchPanic)

	return s
}

// Register associates the given services with the router.
func (s *Server) Register(services ...Service) {
	for _, service := range services {
		service.Register(s.router)
		s.logger.Info(append([]interface{}{"msg", "registered"}, service.Fields()...))
	}
}

// Serve accepts incoming connections, creating a new service goroutine for each. The
// service goroutine reads requests and then call the handler to reply to them.
func (s *Server) Serve() error {
	s.logger.Info("msg", "serving")
	defer s.logger.Info("msg", "serving complete")

	var err error
	if s.listener, err = net.Listen("tcp", s.addr); err != nil {
		return errors.Wrapf(err, "Failed to listen on %s", s.addr)
	}
	s.logger.Info("msg", "listening")

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

// ListenAndServeMetrics runs an HTTP server with /metrics endpoints (e.g. pprof).
func ListenAndServeMetrics() error {
	h := http.NewServeMux()
	h.Handle("/metrics", promhttp.Handler())
	return http.ListenAndServe(":6060", h)
}

// Shutdown stops the http server gracefully without interrupting any active connections.
func (s *Server) Shutdown() {
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(s.shutdownTimeout))
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Error("msg", "Failed to shutdown", "error", err.Error())
	}
	s.logger.Info("msg", "shutdown")
}

// requestLogging is middleware to create a request specific logger, passing
// in through the request's context, as well as log the incoming request.
func (s *Server) requestLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fields := []interface{}{
			"method", r.Method,
			"url", r.URL.String(),
		}
		if len(r.Host) != 0 {
			fields = append(fields, r.Host)
		}
		if len(r.Header.Get("User-Agent")) != 0 {
			fields = append(fields, "user-agent", r.Header.Get("User-Agent"))
		}

		l := s.logger.With(fields...)
		req := r.Clone(log.NewContextWithLogger(r.Context(), l))

		l.Debug("msg", "start")
		next.ServeHTTP(w, req)
		l.Debug("msg", "stop")
	})
}

// requestMetrics is middleware to update the generic http metrics for
// each incoming request.
func (s *Server) requestMetrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Obtain path template & start time of request.
		t := time.Now()
		tmpl := requestPathTemplate(r)

		// Delegate to next handler in middleware chain.
		next.ServeHTTP(w, r)

		// Increment the count and track total time.
		if s.requestCount != nil {
			s.requestCount.WithLabelValues(r.Method, tmpl).Inc()
		}
		if s.requestSeconds != nil {
			s.requestSeconds.WithLabelValues(r.Method, tmpl).Add(float64(time.Since(t).Seconds()))
		}
	})
}

// recoverPanic is middleware for recovering and reporting panics.
func (s *Server) catchPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Install the recovery and reporting function.
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				e, ok := err.(error)
				if ok {
					s.logger.Error("msg", "panic", "error", e.Error())
				} else {
					s.logger.Error("msg", "panic")
				}
			}
		}()
		// Delegate to next handler in middleware chain.
		next.ServeHTTP(w, r)
	})
}

// requestPathTemplate returns the route path template for the given request.
func requestPathTemplate(r *http.Request) string {
	route := mux.CurrentRoute(r)
	if route == nil {
		return ""
	}
	tmpl, _ := route.GetPathTemplate()
	return tmpl
}
