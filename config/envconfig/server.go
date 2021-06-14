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

package envconfig // import "arcadium.dev/core/config/envconfig

import (
	"arcadium.dev/core/config"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const (
	serverPrefix = "server"
)

type (
	// Server holds the configuration settings for a server.
	Server struct {
		addr   string // <PREFIX_>SERVER_ADDR
		cert   string // <PREFIX_>SERVER_CERT
		key    string // <PREFIX_>SERVER_KEY
		cacert string // <PREFIX_>SERVER_CACERT
	}
)

// NewServer returns the server configuration.
func NewServer(opts ...config.Option) (*Server, error) {
	o := &config.Options{}
	for _, opt := range opts {
		opt.Apply(o)
	}
	prefix := o.Prefix() + serverPrefix

	config := struct {
		Addr   string
		Cert   string
		Key    string
		CACert string
	}{}
	if err := envconfig.Process(prefix, &config); err != nil {
		return nil, errors.Wrapf(err, "failed to load %s configuration", prefix)
	}

	return &Server{
		addr:   config.Addr,
		cert:   config.Cert,
		key:    config.Key,
		cacert: config.CACert,
	}, nil
}

// Addr returns the network address the server will listen on.
func (s *Server) Addr() string {
	return s.addr
}

// Cert returns the filename of the certificate.
func (s *Server) Cert() string {
	return s.cert
}

// Key returns the filename of the certificate key.
func (s *Server) Key() string {
	return s.key
}

// CACert returns the filename of the certificate of the client CA.
// This is used when a mutual TLS connection is desired.
func (s *Server) CACert() string {
	return s.cacert
}
