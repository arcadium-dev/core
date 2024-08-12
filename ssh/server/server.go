//Copyright 2023-2024 arcadium.dev <info@arcadium.dev>
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

package server // import "arcadium.dev/core/ssh/server"

import (
	"errors"
	"fmt"
	"net"
	"sync"

	"github.com/rs/zerolog"
	"golang.org/x/crypto/ssh"
)

const (
	defaultAddr         = ":22"
	defaultMaxAuthTries = 3
)

var (
	defaultLogger = zerolog.Nop()
)

type (
	// Server provides an ssh server.
	Server struct {
		mu       sync.RWMutex
		services map[ChannelType]ChannelHandler

		addr   string
		config ssh.ServerConfig
		logger *zerolog.Logger

		listener net.Listener
	}
)

// New creates a new ssh server based on the given options and starts the ssh
// server, listening for ssh connections.
func New(opts ...Option) (*Server, error) {
	s := &Server{
		addr: defaultAddr,
		config: ssh.ServerConfig{
			NoClientAuth: false,
			MaxAuthTries: defaultMaxAuthTries,
		},
		logger: &defaultLogger,
	}
	for _, opt := range opts {
		opt.apply(s)
	}

	listener, err := net.Listen("tcp", s.addr)
	if err != nil {
		return nil, fmt.Errorf("ssh server failed to listen on address '%s', %w", s.addr, err)
	}
	s.listener = listener

	s.logger.Info().Msgf("ssh server created, listening on %s", s.addr)

	return s, nil
}

// Register associates a channel handler with the ssh server.
func (s *Server) Register(services ...ChannelHandler) {
	s.mu.Lock()
	defer s.mu.Unlock()

	for _, service := range services {
		if _, ok := s.services[service.Type()]; !ok {
			s.logger.Warn().Msgf("replacing channel handler for %s", service.Type())
		}
		s.services[service.Type()] = service
	}
}

// Serve begins the process of accepting new ssh connections.
func (s *Server) Serve() error {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			return err
		}
		s.logger.Info().Msgf("ssh connection accepted: %s", conn.RemoteAddr())

		// Run this in a go routine since it could potentially block, for example,
		// during password authentication.
		go s.HandleConn(conn)
	}
	return nil
}

// HandleConn handles an incoming network connection.
func (s *Server) HandleConn(conn net.Conn) {
	defer conn.Close()

	config := s.config
	sshConn, chans, reqs, err := ssh.NewServerConn(conn, &config)
	if err != nil {
		return
	}

	go ssh.DiscardRequests(reqs)

	s.mu.RLock()
	defer s.mu.RUnlock()

	for ch := range chans {
		service, ok := s.services[ChannelType(ch.ChannelType())]
		if !ok {
			ch.Reject(ssh.UnknownChannelType, "unsupported channel type")
			continue
		}
		service.Handle(sshConn, ch)
	}
}

// Close immediately shuts down the ssh server.
func (s *Server) Close() error {
	if err := s.listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
		s.logger.Err(err).Msg("failed to close listener")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, service := range s.services {
		if err := service.Close(); err != nil {
			s.logger.Err(err).Msgf("failed to close channel handler %s", service.Type())
		}
	}

	return nil
}

// Shutdown gracefully shuts down the ssh server.
func (s *Server) Shutdown() error {
	if err := s.listener.Close(); err != nil && !errors.Is(err, net.ErrClosed) {
		s.logger.Err(err).Msg("failed to close listener")
	}

	s.mu.RLock()
	defer s.mu.RUnlock()
	for _, service := range s.services {
		if err := service.Shutdown(); err != nil {
			s.logger.Err(err).Msgf("failed to shutdown channel handler %s", service.Type())
		}
	}

	return nil
}
