//  Copyright 2026 arcadium.dev <info@arcadium.dev>
//
//  Licensed under the Apache License, Version 2.0 (the "License");
//  you may not use this file except in compliance with the License.
//  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package telnet

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"runtime/debug"
	"sync"

	"github.com/globalcyberalliance/telnet-go"
	"github.com/rs/zerolog"
)

const (
	// DefaultAddr defines the default address of the telnet server when one is
	// not provided via the options. This default to the telnet port on the local
	// host.
	DefaultAddr = ":23"
)

type (
	// Session is an alias for telnet-go.Session. This allows for this package to
	// be used without importing telnet-go as well.
	Session = telnet.Session

	// A Handler responds to telnet activity.
	Handler interface {
		ServeTELNET(*Session)
	}

	// HandlerFunc is an adapter to allow the use of an ordinary function as
	// a telnet handler.
	HandlerFunc func(*Session)
)

// ServeTELNET calls f(s).
func (f HandlerFunc) ServeTELNET(s *Session) {
	f(s)
}

type (
	// MiddlewareFunc is a function which takes a Handler and returns a Handler.
	// Typically the returned handler is a closure which does something with the
	// Session passed to it, and then calls the handler passed as the parameter
	// to the MiddlewareFunc.
	MiddlewareFunc = func(Handler) Handler

	// Server represent a telnet server.
	Server struct {
		listenerMutex sync.Mutex
		listener      net.Listener

		serverMutex sync.Mutex
		server      *telnet.Server

		addr       string
		logger     *zerolog.Logger
		middleware []MiddlewareFunc
		service    Service
	}

	// Service defines the methods required by the server to associate the
	// service with the server.
	Service interface {
		// Name returns the name of the service.
		Name() string

		// ServeTELNET provides the telnet handler.
		ServeTELNET(*Session)

		// Shutdown allows the service to stop any long running background processes
		// it may have.
		Shutdown(context.Context)
	}
)

// NewServer create a telnet server with the given server options.
func NewServer(opts ...ServerOption) *Server {
	nop := zerolog.Nop()

	s := &Server{
		addr:   DefaultAddr,
		logger: &nop,
	}

	for _, opt := range opts {
		opt.Apply(s)
	}

	s.server = &telnet.Server{
		Handler:     s.handle,
		ConnContext: s.connContext,
	}
	// It's not really optional.
	s.server.SetLogger(slog.New(slog.DiscardHandler))

	s.logger.Info().Str("address", s.addr).Msg("telnet server created")

	return s
}

// Middleware installs the given middleware with the server.
func (s *Server) Middleware(middleware ...MiddlewareFunc) {
	s.middleware = append(s.middleware, middleware...)
}

// Register associates the given service with the server.
func (s *Server) Register(service Service) {
	s.service = service
	s.logger.Info().Str("service", service.Name()).Msg("telnet service registered")
}

// Serve creates the underlying network connection and starts the telnet
// server.
func (s *Server) Serve() error {
	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listenerMutex.Lock()
	s.listener = listener
	s.listenerMutex.Unlock()

	s.logger.Info().Str("address", s.addr).Str("service", s.service.Name()).Msg("begin serving telnet")
	defer s.logger.Info().Str("address", s.addr).Str("service", s.service.Name()).Msg("serving telnet complete")

	s.serverMutex.Lock()
	defer s.serverMutex.Unlock()

	defer func() {
		if r := recover(); r != nil {
			s.logger.Info().Msg("stacktrace from panic: \n" + string(debug.Stack()))
		}
	}()

	return s.server.Serve(listener)
}

// Shutdown stops the telnet server.
func (s *Server) Shutdown(ctx context.Context) {
	s.service.Shutdown(ctx)

	// This is ugly, but this forces s.server.Serve() to exit.
	s.listenerMutex.Lock()
	if s.listener != nil {
		err := s.listener.Close()
		if err != nil && !errors.Is(err, net.ErrClosed) {
			s.logger.Err(err).Msg("failed to close listener")
		}
	}
	s.listenerMutex.Unlock()

	// This cleans up any outstanding go routines.
	s.serverMutex.Lock()
	defer s.serverMutex.Unlock()
	if err := s.server.Shutdown(); err != nil {
		if !errors.Is(err, net.ErrClosed) {
			s.logger.Err(err).Msg("failed to shutdown telnet server")
		}
	}

	s.logger.Info().Msg("telnet server shutdown")
}

// handle provides the session handler, building a middleware chain
// and then calling the service's handler.
func (s *Server) handle(session *Session) {
	// If the service hasn't been register, log an error.
	if s.service == nil {
		s.logger.Error().Msg("telnet service not registered")
		return
	}

	// Build a chain of functions, starting with the registered service
	// as the last link in the chain, and adding each middleware function in
	// reverst order to the chain.
	var chain Handler = s.service

	// Starting at the end of the middleware slice and working backwards, link
	// the functions together.
	for i := len(s.middleware) - 1; i >= 0; i-- {
		chain = s.middleware[i](chain)
	}

	// Call the chain.
	chain.ServeTELNET(session)
}

// connContext provides a context with an embedded logger containing the
// remote address of the connection.
func (s *Server) connContext(ctx context.Context, conn net.Conn) context.Context {
	logger := s.logger.With().
		Str("remote", conn.RemoteAddr().String()).
		Logger()

	return logger.WithContext(ctx)
}

// Name returns the name of the server.
func (s *Server) Name() string {
	return "telnet server"
}
