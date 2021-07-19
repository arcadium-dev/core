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

package server

import (
	"arcadium.dev/core/config"
	"arcadium.dev/core/grpc"
	"arcadium.dev/core/http"
	"arcadium.dev/core/log"
	"arcadium.dev/core/sql"
)

type (
	// Config implements a server config.
	Config struct {
		logger     log.Config
		db         sql.Config
		grpcServer grpc.Config
		httpServer http.Config
	}

	// ConfigConstructors injects a set of constructors to be used when creating a new server config.
	ConfigConstructors struct {
		// NewLogger is used to creates a new logger config.
		NewLoggerConfig func(...config.Option) (log.Config, error)

		// NewDB is used to creates a new db config.
		NewDBConfig func(...config.Option) (sql.Config, error)

		// NewGRPCServer is used to create a new gRPC server config.
		NewGRPCServerConfig func(...config.Option) (grpc.Config, error)

		// NewHTTPServer is used to create a new HTTP server config.
		NewHTTPServerConfig func(...config.Option) (http.Config, error)
	}
)

// NewConfig returns a new configuration for a server built with the given constructors.
func NewConfig(ctors ConfigConstructors, opts ...config.Option) (*Config, error) {
	l, err := ctors.NewLoggerConfig(opts...)
	if err != nil {
		return nil, err
	}

	d, err := ctors.NewDBConfig(opts...)
	if err != nil {
		return nil, err
	}

	g, err := ctors.NewGRPCServerConfig(opts...)
	if err != nil {
		return nil, err
	}

	h, err := ctors.NewHTTPServerConfig(opts...)
	if err != nil {
		return nil, err
	}

	return &Config{
		logger:     l,
		db:         d,
		grpcServer: g,
		httpServer: h,
	}, nil
}

func (cfg *Config) Logger() log.Config {
	return cfg.logger
}

func (cfg *Config) DB() sql.Config {
	return cfg.db
}

func (cfg *Config) GRPCServer() grpc.Config {
	return cfg.grpcServer
}

func (cfg *Config) HTTPServer() http.Config {
	return cfg.httpServer
}
