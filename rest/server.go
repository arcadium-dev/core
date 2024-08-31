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

// Package rest provides the basis of a restful http service.
package rest // import "arcadium.dev/core/rest"

import (
	"context"
	"fmt"
	"io"
	l "log"
	"os"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"github.com/rs/zerolog"

	"arcadium.dev/core/build"
	"arcadium.dev/core/http/middleware"
	"arcadium.dev/core/http/server"
	"arcadium.dev/core/http/services"
	"arcadium.dev/core/log"
)

type (
	// Server represents the restful api server.
	Server struct {
		Stdout, Stderr io.Writer    // Provides a way for unit tests to capture output to standard file descriptors.
		C              Constructors // Provides a way for unit tests to inject different object constructors.

		interrupt chan os.Signal
		wg        *sync.WaitGroup // To ensure stop isn't called before Start is ready.
		ctx       context.Context

		info   build.Information
		cfg    Config
		logger *zerolog.Logger
		mw     []mux.MiddlewareFunc
	}

	// Constructors provide a way to inject different functions to create server components.
	Constructors struct {
		NewConfig     func(...string) (Config, error)
		NewLogger     func(Config) (*zerolog.Logger, error)
		NewHTTPServer func(context.Context, Config) (*server.Server, error)
	}
)

// NewServer returns a new restful api server.
func NewServer(version, branch, commit, date string, mw ...mux.MiddlewareFunc) *Server {
	name := filepath.Base(os.Args[0])

	s := &Server{
		interrupt: make(chan os.Signal, 1),
		wg:        &sync.WaitGroup{},
		Stdout:    os.Stdout,
		Stderr:    os.Stderr,
		info:      build.Info(name, version, branch, commit, date),
		mw:        mw,
	}

	s.C.NewConfig = NewConfig

	s.C.NewLogger = func(cfg Config) (*zerolog.Logger, error) {
		return log.New(
			log.AsDefault(),
			log.WithLevel(log.ToLevel(cfg.LogLevel())),
			log.WithOutput(s.Stdout),
		)
	}

	s.C.NewHTTPServer = func(ctx context.Context, cfg Config) (*server.Server, error) {
		// Gather the server options.
		var opts []server.Option
		opts = append(opts,
			server.WithAddr(cfg.ServerAddr()),
		)

		// Setup TLS.
		if cfg.TLSCert() != "" && cfg.TLSKey() != "" {
			opts = append(opts, server.WithTLSCert(cfg.TLSCert(), cfg.TLSKey()))
			if cfg.TLSCACert() != "" {
				opts = append(opts, server.WithTLSClientCACert(cfg.TLSCACert()))
			}
			opts = append(opts, server.WithMTLSEnabled(cfg.MTLSEnabled()))
		}

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
			opts = append(opts, server.WithCORSOptions(corsOpts))
		}

		return server.New(ctx, opts...)
	}

	s.wg.Add(1)
	return s
}

// Init initializes the server object.
func (s *Server) Init(prefix ...string) error {
	var (
		err    error
		cancel context.CancelFunc
	)

	// Setup signal handler. Waits for both SIGINT and SIGTERM.
	// SIGTERM is sent by docker to terminate a container.
	s.ctx, cancel = context.WithCancel(context.Background())
	signal.Notify(s.interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() { <-s.interrupt; cancel() }()

	// Setup a temporary logger.
	lg := l.Default()
	lg.SetOutput(s.Stdout)

	// Load the config.
	s.cfg, err = s.C.NewConfig(prefix...)
	if err != nil {
		lg.Printf("error: failed to load config: %s", err)
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create a logger.
	s.logger, err = s.C.NewLogger(s.cfg)
	if err != nil {
		lg.Printf("error: failed to create logger: %s", err)
		return fmt.Errorf("failed to create logger: %w", err)
	}
	s.ctx = s.logger.WithContext(s.ctx)

	s.logger.Info().Msgf("starting %s", s.info)

	return nil
}

func (s Server) Start(extServices ...server.Service) error {
	// Setup http services.
	svcs := []server.Service{
		services.Health{Start: time.Now(), Info: s.info},
		services.Metrics{},
	}
	if s.cfg.PProfEnabled() {
		svcs = append(svcs, services.PProf{})
	}
	if len(extServices) > 0 {
		svcs = append(svcs, extServices...)
	}

	// Create the http server.
	server, err := s.C.NewHTTPServer(s.ctx, s.cfg)
	if err != nil {
		s.logger.Err(err).Msg("failed to create http server")
		return fmt.Errorf("failed to create http server: %w", err)
	}
	mw := []mux.MiddlewareFunc{
		middleware.Recover{Logger: s.logger}.Panics,
		middleware.Logging{Logger: s.logger}.Requests,
		middleware.Metrics,
	}
	if len(s.mw) > 0 {
		mw = append(mw, s.mw...)
	}
	server.Middleware(mw...)

	server.Register(s.ctx, svcs...)

	// Serve.
	result := make(chan error, 1)
	go func() {
		s.wg.Done()
		result <- server.Serve(s.ctx)
	}()

	select {
	// Wait for an interrupt.
	case <-s.ctx.Done():
		server.Shutdown(s.ctx)

	// If the server failed to start,
	case err = <-result:
		if err != nil {
			s.logger.Err(err).Msg("failed to start http server")
		}
	}

	// Shutdown the services.
	for _, svc := range svcs {
		svc.Shutdown(s.ctx)
	}

	return err
}

// Stop halts the server.
func (s Server) Stop() {
	s.wg.Wait()
	close(s.interrupt)
}

func (s Server) Ctx() context.Context { return s.ctx }
