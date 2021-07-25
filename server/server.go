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

package server

//go:generate mockgen -package mockserver -destination ./mock/server.go . GRPCServer,HTTPServer

import (
	"context"
	nethttp "net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"arcadium.dev/core/build"
	"arcadium.dev/core/config"
	"arcadium.dev/core/grpc"
	"arcadium.dev/core/http"
	"arcadium.dev/core/log"
	"arcadium.dev/core/sql"
)

type (
	// Server is a combination grpc and http server.
	Server struct {
		config *Config

		logger log.Logger
		db     sql.DB

		grpcServer GRPCServer
		httpServer HTTPServer

		stop chan struct{}
	}

	// ServerConstructors injects a set of constructors to be used when creating a new server.
	ServerConstructors struct {
		// NewConfig is used to create a new server config.
		NewConfig func(...config.Option) (*Config, error)

		// NewLogger is used to create a new logger.
		NewLogger func(log.Config) log.Logger

		// NewDB is used to create a new DB.
		NewDB func(sql.Config, ...sql.Option) (sql.DB, error)

		// NewGRPCServer is used to create a new gRPC server.
		NewGRPCServer func(grpc.Config, ...grpc.Option) (GRPCServer, error)

		// NewHTTPServer is used to create a new HTTP server.
		NewHTTPServer func(http.Config, ...http.Option) (HTTPServer, error)
	}

	// GRPCServer abstracts a gRPC server.
	GRPCServer interface {
		// Serve is the entry point for a service and will be run it a goroutine.
		// It is passed a channel to communicate the result.
		Serve(result chan<- error)

		// Stop stops the server.
		Stop()

		// Registers the given services with this server.
		Register([]grpc.Service)
	}

	HTTPServer interface {
		// Serve is the entry point for a service and will be run it a goroutine.
		// It is passed a channel to communicate the result.
		Serve(result chan<- error)

		// Stop stops the server.
		Stop()

		// Handle associates the given handler with the server.
		Handle(mux nethttp.Handler)
	}
)

// NewServer returns a combinator grpc and http server.
func New(info build.Information, c ServerConstructors) (*Server, error) {
	var err error

	cfg, err := c.NewConfig()
	if err != nil {
		return nil, err
	}

	logger := c.NewLogger(cfg.Logger())

	db, err := c.NewDB(cfg.DB())
	if err != nil {
		return nil, err
	}

	grpcServer, err := c.NewGRPCServer(cfg.GRPCServer(), grpc.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	httpServer, err := c.NewHTTPServer(cfg.HTTPServer(), http.WithLogger(logger))
	if err != nil {
		return nil, err
	}

	logger.WithFields(info.Fields()).Infof("starting")

	return &Server{
		config:     cfg,
		logger:     logger,
		db:         db,
		grpcServer: grpcServer,
		httpServer: httpServer,
		stop:       make(chan struct{}),
	}, nil
}

// Serve starts the server and will catch os signals. If an os signal is
// caught, the server will be cancelled via it's the context.
func (s *Server) Serve(ctx context.Context) error {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	result := make(chan error, 1)

	go func(result chan<- error) {
		var grpcErr error
		grpcResult := make(chan error)
		go s.grpcServer.Serve(grpcResult)

		var httpErr error
		httpResult := make(chan error)
		go s.httpServer.Serve(httpResult)

		select {
		case <-s.stop:
			s.grpcServer.Stop()
			s.httpServer.Stop()
			grpcErr = <-grpcResult
			httpErr = <-httpResult

		case grpcErr = <-grpcResult:
			s.httpServer.Stop()
			httpErr = <-httpResult

		case httpErr = <-httpResult:
			s.grpcServer.Stop()
			grpcErr = <-grpcResult
		}

		s.db.Close() // Close the db after all server processing has completed.

		switch {
		case grpcErr != nil && httpErr != nil:
			result <- NewErrors(grpcErr, httpErr)
		case grpcErr != nil:
			result <- grpcErr
		case httpErr != nil:
			result <- httpErr
		}
		close(result)
	}(result)

	var err error
	select {
	case <-ctx.Done():
		s.logger.Infof("\n\nContext canceled\n\n")
		s.Stop()
		err = <-result
	case sg := <-sig:
		s.logger.Infof("\n\nSignal received: %s\n\n", sg)
		s.Stop()
		err = <-result
	case err = <-result:
		if err != nil {
			s.logger.Errorf("\n\nError: %s\n\n", err.Error())
		}
	}

	time.Sleep(1 * time.Second) // Give the logs some time to flush.

	return err
}

func (s *Server) Stop() {
	close(s.stop)
}

func (s *Server) Register(services []grpc.Service) {
	s.grpcServer.Register(services)
}

func (s *Server) Handle(mux nethttp.Handler) {
	s.httpServer.Handle(mux)
}

func (s *Server) DB() sql.DB {
	return s.db
}

//---------------------------------------------------------------------------

type Errors struct {
	errs []error
}

func NewErrors(errs ...error) Errors {
	return Errors{errs: errs}
}

func (e Errors) Error() string {
	msg := "Errors:"
	for _, err := range e.errs {
		msg += "\n\t" + err.Error()
	}
	return msg
}
