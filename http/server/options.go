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

package server // import "arcadium.dev/core/http/server"

import (
	"crypto/tls"
	"time"

	"github.com/rs/cors"
)

type (
	// Option provides options for configuring the creation of a http server.
	Option interface {
		apply(*Server)
	}
)

// WithAddr will configure the server with the listen address.
func WithAddr(addr string) Option {
	return newOption(func(s *Server) {
		s.addr = addr
	})
}

// WithTLS will configure the server to require TLS.
func WithTLS(cfg *tls.Config) Option {
	return newOption(func(s *Server) {
		s.server.TLSConfig = cfg
	})
}

// WithCORS will configure the server with the CORS options.
func WithCORS(c *cors.Options) Option {
	return newOption(func(s *Server) {
		s.corsOptions = c
	})
}

// WithReadTimeout sets the http server read timeout.
func WithReadTimeout(timeout time.Duration) Option {
	return newOption(func(s *Server) {
		s.server.ReadTimeout = timeout
	})
}

// WithWriteTimeout sets the http server read timeout.
func WithWriteTimeout(timeout time.Duration) Option {
	return newOption(func(s *Server) {
		s.server.WriteTimeout = timeout
	})
}

// WithShutdownTimeout sets the timout for shutting down the server.
func WithShutdownTimeout(timeout time.Duration) Option {
	return newOption(func(s *Server) {
		s.shutdownTimeout = timeout
	})
}

type (
	option struct {
		f func(*Server)
	}
)

func newOption(f func(*Server)) option {
	return option{f: f}
}

func (o option) apply(s *Server) {
	o.f(s)
}
