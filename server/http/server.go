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
	"net"
	"net/http"
	"time"

	"github.com/pkg/errors"

	"arcadium.dev/core/log"
	"arcadium.dev/core/server"
)

type (

	// Server is an HTTP server to serve HTTP requests.
	Server struct {
		addr   string
		logger log.Logger

		listener net.Listener
		server   *http.Server
	}
)

// New creates an HTTP server with a default handler and has not started to accept requests yet.
func New(config server.Config, opts ...Option) (*Server, error) {
	s := &Server{
		addr:   config.Addr(),
		logger: log.NewNullLogger(),
		server: &http.Server{},
	}

	for _, opt := range opts {
		opt.apply(s)
	}

	// Set up the logging fields.
	fields := map[string]interface{}{
		"server": "http",
		"addr":   s.addr,
	}
	if s.server.TLSConfig != nil {
		fields["tls"] = "enabled"
	}

	// TODO: Do we want to install middleware for metrics, logging or tracing?

	s.logger = s.logger.WithFields(fields)

	return s, nil
}

// Serve accepts incoming connections, creating a new service goroutine for each. The
// service goroutine reads requests and then call the handler to reply to them.
func (s *Server) Serve(result chan<- error) {
	s.logger.Info("serving")
	defer s.logger.Info("serving complete")

	var err error
	if s.listener, err = net.Listen("tcp", s.addr); err != nil {
		result <- errors.Wrapf(err, "Failed to listen on %s", s.addr)
		return
	}
	s.logger.Info("listening")

	if s.server.TLSConfig != nil {
		err = s.server.ServeTLS(s.listener, "", "")
	} else {
		err = s.server.Serve(s.listener)
	}

	if err == http.ErrServerClosed {
		err = nil
	}
	result <- err
}

// Stop shuts down the http server gracefully without interrupting an active connections.
func (s *Server) Stop() {
	// TODO: Do we want to make the shutdown timeout configurable?
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

	if err := s.server.Shutdown(ctx); err != nil {
		s.logger.Errorf("Failed to shutdown http server: %s", err.Error())
	}

	s.logger.Info("stopped")
}

// Handle associates the given handler with the server.
func (s *Server) Handle(mux http.Handler) {
	// Don't overwrite the existing handler if mux is nil.
	if mux == nil {
		return
	}
	s.server.Handler = mux
}
