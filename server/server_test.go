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
	"context"
	"errors"
	"fmt"
	nethttp "net/http"
	"testing"
	"time"

	"github.com/golang/mock/gomock"

	mockgrpc "arcadium.dev/core/grpc/mock"
	mockhttp "arcadium.dev/core/http/mock"
	mocklog "arcadium.dev/core/log/mock"
	mocksql "arcadium.dev/core/sql/mock"

	"arcadium.dev/core/build"
	"arcadium.dev/core/config"
	"arcadium.dev/core/grpc"
	"arcadium.dev/core/http"
	"arcadium.dev/core/log"
	"arcadium.dev/core/sql"
)

//---- New ----

func TestNewSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockLogger := mocklog.NewMockLogger(ctrl)
	mockDB := mocksql.NewMockDB(ctrl)
	grpcServer := &grpc.Server{}
	httpServer := &http.Server{}

	mockLogConfig := mocklog.NewMockConfig(ctrl)
	mockDBConfig := mocksql.NewMockConfig(ctrl)
	mockGRPCServerConfig := mockgrpc.NewMockConfig(ctrl)
	mockHTTPServerConfig := mockhttp.NewMockConfig(ctrl)

	ctor := ServerConstructors{
		NewConfig: func(...config.Option) (*Config, error) {
			return NewConfig(ConfigConstructors{
				NewLoggerConfig:     func(opts ...config.Option) (log.Config, error) { return mockLogConfig, nil },
				NewDBConfig:         func(opts ...config.Option) (sql.Config, error) { return mockDBConfig, nil },
				NewGRPCServerConfig: func(opts ...config.Option) (grpc.Config, error) { return mockGRPCServerConfig, nil },
				NewHTTPServerConfig: func(opts ...config.Option) (http.Config, error) { return mockHTTPServerConfig, nil },
			})
		},
		NewLogger:     func(log.Config) log.Logger { return mockLogger },
		NewDB:         func(sql.Config, ...sql.Option) (sql.DB, error) { return mockDB, nil },
		NewGRPCServer: func(grpc.Config, ...grpc.Option) (GRPCServer, error) { return grpcServer, nil },
		NewHTTPServer: func(http.Config, ...http.Option) (HTTPServer, error) { return httpServer, nil },
	}

	mockLogger.EXPECT().WithFields(gomock.Any()).Return(mockLogger)
	mockLogger.EXPECT().Infof(gomock.Any(), gomock.Any()).Return()

	s, err := New(build.Information{}, ctor)

	if s == nil || err != nil {
		t.Errorf("Unexpected result: %+v %+v", s, err)
	}
	if s.logger != mockLogger {
		t.Errorf("Unexpected result: %+v %+v", s.logger, mockLogger)
	}
	if s.DB() != mockDB {
		t.Errorf("Unexpected result: %+v %+v", s.DB(), mockDB)
	}
	if s.grpcServer != grpcServer {
		t.Errorf("Unexpected result: %+v %+v", s.DB(), mockDB)
	}
	if s.httpServer != httpServer {
		t.Errorf("Unexpected result: %+v %+v", s.DB(), mockDB)
	}
}

func TestNewFailures(t *testing.T) {
	t.Parallel()

	t.Run("Config failure", func(t *testing.T) {
		t.Parallel()

		expectedErr := "NewConfig failure"

		s, err := createServer(t,
			func(...config.Option) (*Config, error) { return nil, errors.New(expectedErr) },
		)
		if s != nil || err.Error() != expectedErr {
			t.Errorf("Error expected: %+v %+v", s, err)
		}
	})

	t.Run("DB failure", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		ctrl.Finish()
		mockLogConfig := mocklog.NewMockConfig(ctrl)
		expectedErr := "NewDB failure"

		s, err := createServer(t,
			func(...config.Option) (*Config, error) {
				return createConfig(t, func(opts ...config.Option) (log.Config, error) { return mockLogConfig, nil })
			},
			func(sql.Config, ...sql.Option) (sql.DB, error) {
				return nil, errors.New(expectedErr)
			},
		)
		if s != nil || err.Error() != expectedErr {
			t.Errorf("Error expected: %+v %+v", s, err)
		}
	})

	t.Run("HTTP server failure", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		ctrl.Finish()
		mockLogConfig := mocklog.NewMockConfig(ctrl)
		expectedErr := "NewHTTPServer failure"

		s, err := createServer(t,
			func(...config.Option) (*Config, error) {
				return createConfig(t, func(opts ...config.Option) (log.Config, error) { return mockLogConfig, nil })
			},
			func(http.Config, ...http.Option) (HTTPServer, error) {
				return nil, errors.New(expectedErr)
			},
		)
		if s != nil || err.Error() != expectedErr {
			t.Errorf("Error expected: %+v %+v", s, err)
		}
	})

	t.Run("GRPC server failure", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		ctrl.Finish()
		mockLogConfig := mocklog.NewMockConfig(ctrl)
		expectedErr := "NewGRPCServer failure"

		s, err := createServer(t,
			func(...config.Option) (*Config, error) {
				return createConfig(t, func(opts ...config.Option) (log.Config, error) { return mockLogConfig, nil })
			},
			func(grpc.Config, ...grpc.Option) (GRPCServer, error) {
				return nil, errors.New(expectedErr)
			},
		)
		if s != nil || err.Error() != expectedErr {
			t.Errorf("Error expected: %+v %+v", s, err)
		}
	})
}

func createServer(t *testing.T, newFuncs ...interface{}) (*Server, error) {
	t.Helper()

	ctor := ServerConstructors{
		NewConfig:     func(...config.Option) (*Config, error) { return nil, nil },
		NewLogger:     func(log.Config) log.Logger { return nil },
		NewDB:         func(sql.Config, ...sql.Option) (sql.DB, error) { return nil, nil },
		NewGRPCServer: func(grpc.Config, ...grpc.Option) (GRPCServer, error) { return nil, nil },
		NewHTTPServer: func(http.Config, ...http.Option) (HTTPServer, error) { return nil, nil },
	}

	for _, newFunc := range newFuncs {
		switch f := newFunc.(type) {
		case func(...config.Option) (*Config, error):
			ctor.NewConfig = f
		case func(log.Config) log.Logger:
			ctor.NewLogger = f
		case func(sql.Config, ...sql.Option) (sql.DB, error):
			ctor.NewDB = f
		case func(grpc.Config, ...grpc.Option) (GRPCServer, error):
			ctor.NewGRPCServer = f
		case func(http.Config, ...http.Option) (HTTPServer, error):
			ctor.NewHTTPServer = f
		default:
			t.Errorf("Unexpected newFunc")
		}
	}
	return New(build.Information{}, ctor)
}

// ---- Serve ----

func TestServeSuccess(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	s := setupServer(t, ctrl,
		func(config grpc.Config, opts ...grpc.Option) (GRPCServer, error) {
			opts = append(opts, grpc.WithInsecure())
			return grpc.New(config, opts...)
		},
		func(config http.Config, opts ...http.Option) (HTTPServer, error) {
			return http.New(config, opts...)
		},
		true,
	)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(500*time.Millisecond))
	defer cancel()

	err := s.Serve(ctx)

	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

func TestServerError(t *testing.T) {
	t.Parallel()

	t.Run("Stop GRPC error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		expectedErr := "grpc error"

		s := setupServer(t, ctrl,
			func(config grpc.Config, opts ...grpc.Option) (GRPCServer, error) {
				return mockGRPCServer{err: errors.New(expectedErr)}, nil
			},
			func(config http.Config, opts ...http.Option) (HTTPServer, error) {
				return mockHTTPServer{sleep: 500}, nil
			},
			false,
		)

		err := s.Serve(context.Background())
		if err.Error() != expectedErr {
			t.Errorf("Expected %s, actual %s", expectedErr, err.Error())
		}
	})

	t.Run("Stop HTTP error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		expectedErr := "http error"

		s := setupServer(t, ctrl,
			func(config grpc.Config, opts ...grpc.Option) (GRPCServer, error) {
				return mockGRPCServer{sleep: 500}, nil
			},
			func(config http.Config, opts ...http.Option) (HTTPServer, error) {
				return mockHTTPServer{err: errors.New(expectedErr)}, nil
			},
			false,
		)

		err := s.Serve(context.Background())
		if err.Error() != expectedErr {
			t.Errorf("Expected %s, actual %s", expectedErr, err.Error())
		}
	})

	t.Run("Stop GRPC and HTTP errors", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		expectedGRPCErr := "grpc error"
		expectedHTTPErr := "http error"

		s := setupServer(t, ctrl,
			func(config grpc.Config, opts ...grpc.Option) (GRPCServer, error) {
				return mockGRPCServer{err: errors.New(expectedGRPCErr)}, nil
			},
			func(config http.Config, opts ...http.Option) (HTTPServer, error) {
				return mockHTTPServer{err: errors.New(expectedHTTPErr)}, nil
			},
			false,
		)

		err := s.Serve(context.Background())
		if err == nil {
			t.Errorf("Error expected")
		}
		expectedErrs := fmt.Sprintf("Errors:\n\t%s\n\t%s", expectedGRPCErr, expectedHTTPErr)
		if err.Error() != expectedErrs {
			t.Errorf("Expected %s, actual %s", expectedErrs, err.Error())
		}
	})
}

type (
	grpcServerCtor func(config grpc.Config, opts ...grpc.Option) (GRPCServer, error)
	httpServerCtor func(config http.Config, opts ...http.Option) (HTTPServer, error)
)

func setupServer(t *testing.T, ctrl *gomock.Controller, newGRPCServer grpcServerCtor, newHTTPServer httpServerCtor, realServers bool) *Server {
	t.Helper()

	mockDB := mocksql.NewMockDB(ctrl)

	mockLogConfig := mocklog.NewMockConfig(ctrl)
	mockDBConfig := mocksql.NewMockConfig(ctrl)
	mockGRPCServerConfig := mockgrpc.NewMockConfig(ctrl)
	mockHTTPServerConfig := mockhttp.NewMockConfig(ctrl)

	ctor := ServerConstructors{
		NewConfig: func(...config.Option) (*Config, error) {
			return NewConfig(ConfigConstructors{
				NewLoggerConfig:     func(opts ...config.Option) (log.Config, error) { return mockLogConfig, nil },
				NewDBConfig:         func(opts ...config.Option) (sql.Config, error) { return mockDBConfig, nil },
				NewGRPCServerConfig: func(opts ...config.Option) (grpc.Config, error) { return mockGRPCServerConfig, nil },
				NewHTTPServerConfig: func(opts ...config.Option) (http.Config, error) { return mockHTTPServerConfig, nil },
			})
		},
		NewLogger:     func(log.Config) log.Logger { return log.NewNullLogger() },
		NewDB:         func(sql.Config, ...sql.Option) (sql.DB, error) { return mockDB, nil },
		NewGRPCServer: newGRPCServer,
		NewHTTPServer: newHTTPServer,
	}

	if realServers {
		mockGRPCServerConfig.EXPECT().Addr().Return(":4201")
		mockHTTPServerConfig.EXPECT().Addr().Return(":8080")
	}

	mockDB.EXPECT().Close().Return(nil)

	s, err := New(build.Information{}, ctor)
	if s == nil || err != nil {
		t.Errorf("Unexpected result: %+v %+v", s, err)
	}

	return s
}

type mockGRPCServer struct {
	sleep time.Duration
	err   error
}

func (s mockGRPCServer) Serve(result chan<- error) {
	time.Sleep(s.sleep * time.Millisecond)
	if s.err != nil {
		result <- s.err
	}
	close(result)
}

func (s mockGRPCServer) Stop()                   {}
func (s mockGRPCServer) Register([]grpc.Service) {}

type mockHTTPServer struct {
	sleep time.Duration
	err   error
}

func (s mockHTTPServer) Serve(result chan<- error) {
	time.Sleep(s.sleep * time.Millisecond)
	if s.err != nil {
		result <- s.err
	}
	close(result)
}

func (s mockHTTPServer) Stop()                  {}
func (s mockHTTPServer) Handle(nethttp.Handler) {}
