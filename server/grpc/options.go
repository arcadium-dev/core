// Copyright 2021 Ian Cahoon <icahoon@gmail.com>
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

package grpc // import "arcadium.dev/core/server/grpc"

import (
	"arcadium.dev/core/log"
	"arcadium.dev/core/trace"
)

type (
	// Option provides options for configuring the creation of a gRPC server.
	Option interface {
		apply(*Server)
	}

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

// WithoutReflection disables the registration of the reflection service.
func WithoutReflection() Option {
	return newOption(func(s *Server) {
		s.reflection = false
	})
}

// WithInsecure will configure the server to allow insecure connections.
func WithInsecure() Option {
	return newOption(func(s *Server) {
		s.insecure = true
	})
}

// WithMetrics will configure a unary server interceptor which provides server
// side metrics.
func WithMetrics() Option {
	return newOption(func(s *Server) {
		s.metrics = true
	})
}

// WithLogger will configure a unary server interceptor which provides logging
// of each gRPC request, and embeds a logger into the request's context.
func WithLogger(l log.Logger) Option {
	return newOption(func(s *Server) {
		s.logger = l
	})
}

// WithTrace will configure a unary server interceptor which provides tracing
// of each gRPC request, and embeds trace info into the request's context.
func WithTracer(t trace.Tracer) Option {
	return newOption(func(s *Server) {
		s.tracer = t
	})
}
