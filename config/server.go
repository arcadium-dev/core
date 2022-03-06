// Copyright 2021-2022 arcadium.dev <info@arcadium.dev>
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

package config // import "arcadium.dev/core/config

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type (
	// Server holds the configuration settings for a server.
	Server struct {
		addr   string
		cert   string
		key    string
		cacert string
	}
)

const (
	serverPrefix = "server"
)

// NewServer returns the server configuration.
func NewServer(opts ...Option) (Server, error) {
	o := &Options{}
	for _, opt := range opts {
		opt.Apply(o)
	}
	prefix := o.Prefix + serverPrefix

	config := struct {
		Addr   string `required:"true"`
		Cert   string
		Key    string
		CACert string
	}{}
	if err := envconfig.Process(prefix, &config); err != nil {
		return Server{}, fmt.Errorf("failed to load %s configuration: %w", prefix, err)
	}

	return Server{
		addr:   strings.TrimSpace(config.Addr),
		cert:   strings.TrimSpace(config.Cert),
		key:    strings.TrimSpace(config.Key),
		cacert: strings.TrimSpace(config.CACert),
	}, nil
}

// Addr returns the network address the server will listen on. The value is set
// from the <PREFIX_>SERVER_ADDR environment variable.
func (s Server) Addr() string {
	return s.addr
}

// Cert returns the path of the certificate file. The value is set from the
// <PREFIX_>SERVER_CERT environment variable.
func (s Server) Cert() string {
	return s.cert
}

// Key returns the path of the certificate key file. The value is set from the
// <PREFIX_>SERVER_KEY environment variable.
func (s Server) Key() string {
	return s.key
}

// CACert returns the path of the certificate of the CA certificate file. This
// is used when creating a TLS connection with an entity that is presenting a
// certificate that is not signed by a well known CA available in the OS CA
// bundle. The value is set from the <PREFIX_>SERVER_CACERT environment
// variable.
func (s Server) CACert() string {
	return s.cacert
}
