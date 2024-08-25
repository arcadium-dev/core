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
	"io"
	"time"
)

type (
	// Option provides options for configuring the creation of a multiprotocol server.
	Option interface {
		apply(*MultiprotocolServer)
	}

	option struct {
		f func(*MultiprotocolServer)
	}
)

func newOption(f func(*MultiprotocolServer)) option {
	return option{f: f}
}

func (o option) apply(s *MultiprotocolServer) {
	o.f(s)
}

// WithLogLevel ...
func WithLogLevel(logLevel string) Option {
	return newOption(func(s *MultiprotocolServer) {
		s.loglevel = logLevel
	})
}

// WithProtocolServer ...
func WithProtocolServer(pserver ProtocolServer) Option {
	return newOption(func(s *MultiprotocolServer) {
		if pserver != nil {
			s.servers = append(s.servers, pserver)
		}
	})
}

// WithStdout ...
func WithStdout(stdout io.Writer) Option {
	return newOption(func(s *MultiprotocolServer) {
		s.stdout = stdout
	})
}

func WithShutdownTimeout(d time.Duration) Option {
	return newOption(func(s *MultiprotocolServer) {
		s.shutdownTimeout = d
	})
}
