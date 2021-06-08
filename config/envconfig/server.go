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

package envconfig // import "arcadium.dev/core/config/envconfig

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const (
	serverPrefix = "server"
)

type (
	Server struct {
		addr   string // <PREFIX_>SERVER_ADDR
		cert   string // <PREFIX_>SERVER_CERT
		key    string // <PREFIX_>SERVER_KEY
		cacert string // <PREFIX_>SERVER_CACERT
	}
)

func NewServer(opts ...Option) (*Server, error) {
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

	return &Server{
		addr:   config.Addr,
		cert:   config.Cert,
		key:    config.Key,
		cacert: config.CACert,
	}, nil
}

func (s *Server) Addr() string {
	return s.addr
}

func (s *Server) Cert() string {
	return s.cert
}

func (s *Server) Key() string {
	return s.key
}

func (s *Server) CACert() string {
	return s.cacert
}
