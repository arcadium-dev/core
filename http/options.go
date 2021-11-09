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

	"github.com/prometheus/client_golang/prometheus"
)

type (
	// ServerOption provides options for configuring the creation of a http server.
	ServerOption interface {
		apply(*Server)
	}
)

// WithTLS will configure the server to require TLS.
func WithTLS(cfg *tls.Config) ServerOption {
	return newServerOption(func(s *Server) {
		s.server.TLSConfig = cfg
	})
}

// WithShutdownTimeout sets the timout for shutting down the server.
func WithShutdownTimeout(timeout time.Duration) ServerOption {
	return newServerOption(func(s *Server) {
		s.shutdownTimeout = timeout
	})
}

// WithMetrics provides generic request counters to the server.
func WithMetrics(requestCount, requestSeconds *prometheus.CounterVec) ServerOption {
	return newServerOption(func(s *Server) {
		s.requestCount = requestCount
		s.requestSeconds = requestSeconds
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
