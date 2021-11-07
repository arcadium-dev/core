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
	"crypto/tls"
	"time"
)

type (
	// Option provides options for configuring the creation of a http server.
	Option interface {
		apply(*Server)
	}
)

// WithTLS will configure the server to require TLS.
func WithTLS(cfg *tls.Config) Option {
	return newOption(func(s *Server) {
		s.server.TLSConfig = cfg
	})
}

// WithTimeout sets the timout for stopping the server.
func WithTimeout(timeout time.Duration) Option {
	return newOption(func(s *Server) {
		s.timeout = timeout
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
