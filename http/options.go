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

package http // import "arcadium.dev/core/http"

import (
	"crypto/tls"
	"time"

	"arcadium.dev/core/log"
)

type (
	// ServerOption provides options for configuring the creation of a http server.
	ServerOption interface {
		apply(*Server)
	}
)

// WithServerAddr will configure the server with the listen address.
func WithServerAddr(addr string) ServerOption {
	return newServerOption(func(s *Server) {
		s.addr = addr
	})
}

// WithServerTLS will configure the server to require TLS.
func WithServerTLS(cfg *tls.Config) ServerOption {
	return newServerOption(func(s *Server) {
		s.server.TLSConfig = cfg
	})
}

// WithServerShutdownTimeout sets the timout for shutting down the server.
func WithServerShutdownTimeout(timeout time.Duration) ServerOption {
	return newServerOption(func(s *Server) {
		s.shutdownTimeout = timeout
	})
}

// WithServerLogger provides a logger to the server.
func WithServerLogger(logger log.Logger) ServerOption {
	return newServerOption(func(s *Server) {
		s.logger = logger
	})
}

type (
	serverOption struct {
		f func(*Server)
	}
)

func newServerOption(f func(*Server)) serverOption {
	return serverOption{f: f}
}

func (o serverOption) apply(s *Server) {
	o.f(s)
}
