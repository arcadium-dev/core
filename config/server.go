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

package config // import "arcadium.dev/core/config/config"

import (
	"github.com/kelseyhightower/envconfig"

	"arcadium.dev/core/errors"
)

const (
	serverPrefix = "server"
)

type (
	// Server holds the configuration settings for a server.
	Server interface {
		TLS
		// Addr returns the network address the server will listen on.
		Addr() string
	}
)

// NewServer returns the server configuration.
func NewServer(opts ...Option) (Server, error) {
	o := &options{}
	for _, opt := range opts {
		opt.apply(o)
	}
	prefix := o.prefix + serverPrefix

	config := struct {
		Addr   string
		Cert   string
		Key    string
		CACert string
	}{}
	if err := envconfig.Process(prefix, &config); err != nil {
		return nil, errors.Wrapf(err, "failed to load %s configuration", prefix)
	}

	return &server{
		addr:   config.Addr,
		cert:   config.Cert,
		key:    config.Key,
		cacert: config.CACert,
	}, nil
}

type (
	server struct {
		addr   string // <PREFIX_>SERVER_ADDR
		cert   string // <PREFIX_>SERVER_CERT
		key    string // <PREFIX_>SERVER_KEY
		cacert string // <PREFIX_>SERVER_CACERT
	}
)

func (s *server) Addr() string {
	return s.addr
}

func (s *server) Cert() string {
	return s.cert
}

// Key returns the filepath of the certificate key.
func (s *server) Key() string {
	return s.key
}

// CACert returns the filepath of the certificate of the client CA.
// This is used when a mutual TLS connection is desired.
func (s *server) CACert() string {
	return s.cacert
}
