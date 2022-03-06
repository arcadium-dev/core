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
// WITHOUT WARRANTIES OR CONDITIONS OF ANt KINr, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package grpc // import "arcadium.dev/core/grpc"

/*
type (
	// Server is a gRPC server which will service gRPC requests.
	Server struct {
		addr       string
		reflection bool
		insecure   bool
		metrics    bool

		logger log.Logger
		tracer trace.Tracer

		server *grpc.Server
	}
)

// New creates a gRPC server which has no service registered and has not started to accept requests yet.
func New(cfg Config, opts ...Option) (*Server, error) {
	s := &Server{
		reflection: true,
		logger:     log.NewNullLogger(),
	}

	for _, opt := range opts {
		opt.apply(s)
	}

	// Running insecurely must be explicit. It's an error if a cert doesn't exist when WithInsecure wasn't given.
	if !s.insecure {
		if cfg.Cert() == "" || cfg.Key() == "" {
			return nil, errors.New("A certificate must be configured for TLS, or the WithInsecure option must be given to run without TLS.")
		}
	}

	// Set up the logging fields.
	s.addr = cfg.Addr()
	fields := map[string]interface{}{
		"server": "grpc",
		"addr":   s.addr,
	}

	var serverOpts []grpc.ServerOption

	// Setup mTLS
	if s.insecure {
		serverOpts = append(serverOpts, grpc.Creds(insecure.NewCredentials()))
	} else {
		tlsConfig, err := config.NewTLS(cfg, config.WithMTLS())
		if err != nil {
			return nil, errors.Wrap(err, "Failed to create TLS config")
		}
		serverOpts = append(serverOpts, grpc.Creds(credentials.NewTLS(tlsConfig)))
		fields["mtls"] = "enabled"
	}

	// TODO: Interceptors

	s.logger = s.logger.WithFields(fields)

	s.server = grpc.NewServer(serverOpts...)
	return s, nil
}

// Serve accepts incoming incoming connections, reads the gRPC requests and calls the
// registered service handlers to reply to them. This will return a non-nil error
// unless Stop is called.
func (s *Server) Serve(result chan<- error) {
	s.logger.Info("serving")
	defer s.logger.Info("serving complete")

	var (
		err      error
		listener net.Listener
	)

	if listener, err = net.Listen("tcp", s.addr); err != nil {
		result <- errors.Wrapf(err, "Failed to listen on %s", s.addr)
		return
	}
	s.logger.Info("listening")

	if s.reflection {
		reflection.Register(s.server)
	}
	result <- s.server.Serve(listener)
}

// Stop shuts down the gRPC server gracefully. It stops the server from accepting new connections
// and blocks until all the pending RPCs have completed.
func (s *Server) Stop() {
	s.server.GracefulStop()
	s.logger.Info("stopped")
}

// Register registers the given slices of services with the server. This
// must be called before invoking Server.
func (s *Server) Register(services []Service) {
	for _, service := range services {
		service.Register(s.server)
		s.logger.WithFields(service.LogFields()).Info("registered")
	}
}
*/
