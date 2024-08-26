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

package mpserver // import "arcadium.dev/core/mpserver"

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/rs/zerolog"

	"arcadium.dev/core/build"
	"arcadium.dev/core/log"
)

type (
	// MultiprotocolServer manages servers that provide services with mutliple
	// protocols.
	MultiprotocolServer struct {
		// configurable via options
		loglevel        string
		servers         []ProtocolServer
		shutdownTimeout time.Duration
		stdout          io.Writer

		// private data
		ctx       context.Context
		info      build.Information
		interrupt chan os.Signal
		logger    *zerolog.Logger
	}

	// ProtocolServer defines the behavior expended from a protocol server.
	ProtocolServer interface {
		// Serve starts the server. This will be run in its own go routine.
		Serve(context.Context, build.Information) error

		// Shutdown a protocol server. Calling shutdown for a server that returns
		// an erro from Serve must be a noop.
		Shutdown(context.Context)
	}
)

const (
	defaultLogLevel        = "info"
	defaultShutdownTimeout = 10 * time.Second
)

// New returns a new multiprotocol server.
func New(version, branch, commit, date string, opts ...Option) (*MultiprotocolServer, error) {
	s := &MultiprotocolServer{
		loglevel:        defaultLogLevel,
		shutdownTimeout: defaultShutdownTimeout,
		stdout:          os.Stdout,

		info:      build.Info(filepath.Base(os.Args[0]), version, branch, commit, date),
		interrupt: make(chan os.Signal, 1),
	}

	// Load options.
	for _, opt := range opts {
		opt.apply(s)
	}

	var (
		err    error
		cancel context.CancelFunc
	)

	// Setup signal handler. Waits for both SIGINT and SIGTERM.
	// SIGTERM is sent by docker to terminate a container.
	s.ctx, cancel = context.WithCancel(context.Background())
	signal.Notify(s.interrupt, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-s.interrupt
		cancel()
	}()

	// Create a logger.
	if s.logger, err = log.New(log.AsDefault(), log.WithLevel(log.ToLevel(s.loglevel)), log.WithOutput(s.stdout)); err != nil {
		return nil, fmt.Errorf("failed to create logger: %w", err)
	}
	s.ctx = s.logger.WithContext(s.ctx)

	return s, nil
}

// Serve starts the protocol servers, waits for them to start (or fail), and
// then shuts down the servers.
func (s *MultiprotocolServer) Serve() error {
	if len(s.servers) == 0 {
		return fmt.Errorf("exiting, nothing to server")
	}

	s.logger.Info().Msgf("starting %s", s.info)

	result := make(chan error, len(s.servers))
	for _, server := range s.servers {
		go func() {
			result <- server.Serve(s.ctx, s.info)
		}()
	}

	var err error
	select {
	// Wait for an interrupt.
	case <-s.ctx.Done():

	// If a protocol server fails to start.
	case err = <-result:
		if err != nil {
			s.logger.Err(err).Msg("failed to start server")
		}
	}

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(s.shutdownTimeout))
	defer cancel()
	for _, server := range s.servers {
		server.Shutdown(ctx)
	}

	return err
}

// Shutdown halts the server.
func (s MultiprotocolServer) Shutdown() {
	close(s.interrupt)
}

// Ctx returns the context used by this server.
func (s MultiprotocolServer) Ctx() context.Context { return s.ctx }
