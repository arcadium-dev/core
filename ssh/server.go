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
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog"

	"arcadium.dev/core/build"
	"arcadium.dev/core/http/middleware"
	"arcadium.dev/core/http/server"
	"arcadium.dev/core/http/services"
	"arcadium.dev/core/log"
)

type (
	// MultiprotocolServer ...
	MultiprotocolServer struct {
		LogLevel string

		stdout io.Writer // Provides a way for unit tests to capture output to standard file descriptors.

		interrupt chan os.Signal
		ctx       context.Context
		info      build.Information
		logger    *zerolog.Logger
	}

	// Constructors provide a way to inject different functions to create server components.
	Constructors struct {
		NewConfig func(...string) (Config, error)
		NewLogger func(Config) (*zerolog.Logger, error)
	}

	// ProtocolServer ...
	ProtocolServer interface {
		Start(context.Context, build.Information)
		Done() <-chan error
		Shutdown(context.Context)
	}
)

// New returns a new multiprotocol server.
func NewServer(version, branch, commit, date string) *MultiprotocolServer {
	return &MultiprotocolServer{
		interrupt: make(chan os.Signal, 1),
		stdout:    os.Stdout,
		info:      build.Info(filepath.Base(os.Args[0]), version, branch, commit, date),
	}
}

// Init initializes the server object.
func (s *MultiprotocolServer) Init(prefix ...string) error {
	var (
		err    error
		cancel context.CancelFunc
	)

	// Setup signal handler. Waits for both SIGINT and SIGTERM.
	// SIGTERM is sent by docker to terminate a container.
	s.ctx, cancel = context.WithCancel(context.Background())
	signal.Notify(s.interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() { <-s.interrupt; cancel() }()

	// Create a logger.
	if s.logger, err = log.New(log.AsDefault(), log.WithLevel(log.ToLevel(cfg.LogLevel())), log.WithOutput(s.Stdout)); err != nil {
		return fmt.Errorf("failed to create logger: %w", err)
	}
	s.ctx = s.logger.WithContext(s.ctx)

	// Create the http server.
	if s.httpServer, err = s.C.NewHTTPServer(s.ctx, s.cfg); err != nil {
		return fmt.Errorf("failed to create http server: %w", err)
	}

	s.logger.Info().Msgf("starting %s", s.info)

	return nil
}

func (s Server) Start() error {
	// Setup http services.
	svcs := []server.Service{
		services.Health{Start: time.Now(), Info: s.info},
		services.Metrics{},
	}
	if s.cfg.PProfEnabled() {
		svcs = append(svcs, services.PProf{})
	}
	s.server.Register(svcs...)

	// Setup the http middleware.
	mw := []mux.MiddlewareFunc{
		middleware.Recover{Logger: s.logger}.Panics,
		middleware.Logging{Logger: s.logger}.Requests,
		middleware.Metrics,
	}
	server.Middleware(mw...)

	// Create the http server.
	httpServer, err := s.C.NewHTTPServer(s.ctx, s.cfg)
	if err != nil {
		s.logger.Err(err).Msg("failed to create http server")
		return fmt.Errorf("failed to create http server: %w", err)
	}

	// Serve.
	result := make(chan error, 1)
	go func() {
		s.wg.Done()
		result <- httpServer.Serve()
	}()

	select {
	// Wait for an interrupt.
	case <-s.ctx.Done():
		httpServer.Shutdown()

	// If the http server fails to start...
	case err = <-result:
		if err != nil {
			s.logger.Err(err).Msg("failed to start http server")
		}
	}

	// Shutdown the services.
	for _, svc := range svcs {
		svc.Shutdown()
	}

	return err
}

// Stop halts the server.
func (s Server) Stop() {
	s.wg.Wait()
	close(s.interrupt)
}

func (s Server) Ctx() context.Context { return s.ctx }
