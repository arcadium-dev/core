// Copyright 2023-2024 arcadium.dev <info@arcadium.dev>
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
	"github.com/rs/zerolog"
	"golang.org/x/crypto/ssh"
)

type (
	// Option provides options for configuring the creation of an ssh server.
	Option interface {
		apply(*Server)
	}

	// PasswordCallback defines a function used to do authentication via a password.
	PasswordCallback func(conn ssh.ConnMetadata, pass []byte) (*ssh.Permissions, error)

	// PublicKeyCallback defines a function used to do public key authentication.
	PublicKeyCallback func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error)
)

// WithAddr will configure the server with the listen address.
func WithAddr(addr string) Option {
	return newOption(func(s *Server) {
		if addr != "" {
			s.addr = addr
		}
	})
}

// WithHostKeys provides one or more host keys for the server.
func WithHostKeys(keys ...ssh.Signer) Option {
	return newOption(func(s *Server) {
		for _, key := range keys {
			if key == nil {
				panic("nil host key passed as an option to the ssh server")
			}
			s.config.AddHostKey(key)
		}
	})
}

// WithPublicKeyAuthn provides a public key authentication callback function.
func WithPublicKeyAuthn(handler PublicKeyCallback) Option {
	return newOption(func(s *Server) {
		if s.config.PublicKeyCallback != nil {
			panic("public key callback already defined in ssh server config")
		}
		s.config.PublicKeyCallback = handler
	})
}

// WithPasswordAuthn provides a password authentication callback function.
func WithPasswordAuthn(handler PasswordCallback) Option {
	return newOption(func(s *Server) {
		if s.config.PasswordCallback != nil {
			panic("password callback already defined in ssh server config")
		}
		s.config.PasswordCallback = handler
	})
}

// WithServerLogger provides a logger to the server.
func WithLogger(logger *zerolog.Logger) Option {
	return newOption(func(s *Server) {
		s.logger = logger
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
