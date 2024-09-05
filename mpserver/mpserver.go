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
	"sync"
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
		shutdownTimeout time.Duration
		stdout          io.Writer

		// private data
		ctx       context.Context
		info      build.Information
		interrupt chan os.Signal
		logger    *zerolog.Logger

		servers *servers
	}

	// ProtocolServer defines the behavior expended from a protocol server.
	ProtocolServer interface {
		// Serve starts the server. This will be run in its own go routine.
		Serve(context.Context) error

		// Shutdown a protocol server. Calling shutdown for a server that returns
		// an erro from Serve must be a noop.
		Shutdown(context.Context)

		// Name returns the name of the server.
		Name() string
	}

	servers struct {
		mu      sync.RWMutex
		servers []ProtocolServer
	}
)

func (s *servers) append(svrs []ProtocolServer) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.servers = append(s.servers, svrs...)
}

func (s *servers) len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.servers)
}

func (s *servers) get(i int) ProtocolServer {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.servers[i]
}

const (
	DefaultLogLevel        = "info"
	DefaultShutdownTimeout = 10 * time.Second
)

// New returns a new multiprotocol server.
func New(version, branch, commit, date string, opts ...Option) (*MultiprotocolServer, error) {
	s := &MultiprotocolServer{
		loglevel:        DefaultLogLevel,
		shutdownTimeout: DefaultShutdownTimeout,
		stdout:          os.Stdout,
		info:            build.Info(filepath.Base(os.Args[0]), version, branch, commit, date),
		interrupt:       make(chan os.Signal, 1),
		servers:         &servers{},
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

	s.logger.Info().Msgf("starting %s", s.info)

	return s, nil
}

func (s *MultiprotocolServer) Register(ctx context.Context, servers ...ProtocolServer) {
	s.servers.append(servers)
	for _, server := range servers {
		zerolog.Ctx(ctx).Info().Msgf("protocol server registered: %s", server.Name())
	}
}

// Serve starts the protocol servers, waits for them to start (or fail), and
// then shuts down the servers.
func (s MultiprotocolServer) Serve() error {
	if s.servers.len() == 0 {
		return fmt.Errorf("exiting, nothing to server")
	}

	l := s.servers.len()
	result := make(chan error, l)
	for i := 0; i < l; i++ {
		server := s.servers.get(i)
		go func() {
			result <- server.Serve(s.ctx)
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
	for i := 0; i < s.servers.len(); i++ {
		s.servers.get(i).Shutdown(ctx)
	}

	return err
}

// Shutdown halts the server.
func (s MultiprotocolServer) Shutdown() {
	close(s.interrupt)
}

// Ctx returns the context used by this server.
func (s MultiprotocolServer) Ctx() context.Context { return s.ctx }
