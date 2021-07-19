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
	"errors"
	"testing"

	"github.com/golang/mock/gomock"

	mockgrpc "arcadium.dev/core/grpc/mock"
	mockhttp "arcadium.dev/core/http/mock"
	mocklog "arcadium.dev/core/log/mock"
	mocksql "arcadium.dev/core/sql/mock"

	"arcadium.dev/core/config"
	"arcadium.dev/core/grpc"
	"arcadium.dev/core/http"
	"arcadium.dev/core/log"
	"arcadium.dev/core/sql"
)

func TestNewServerConfigSuccess(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)

	mockLoggerConfig := mocklog.NewMockConfig(ctrl)
	mockDBConfig := mocksql.NewMockConfig(ctrl)
	mockGRPCServerConfig := mockgrpc.NewMockConfig(ctrl)
	mockHTTPServerConfig := mockhttp.NewMockConfig(ctrl)

	ctors := ConfigConstructors{
		NewLoggerConfig: func(opts ...config.Option) (log.Config, error) {
			return mockLoggerConfig, nil
		},
		NewDBConfig: func(opts ...config.Option) (sql.Config, error) {
			return mockDBConfig, nil
		},
		NewGRPCServerConfig: func(opts ...config.Option) (grpc.Config, error) {
			return mockGRPCServerConfig, nil
		},
		NewHTTPServerConfig: func(opts ...config.Option) (http.Config, error) {
			return mockHTTPServerConfig, nil
		},
	}

	cfg, err := NewConfig(ctors)
	if cfg == nil || err != nil {
		t.Errorf("unexpected failure: %+v %+v", cfg, err)
	}
	if cfg.Logger() != mockLoggerConfig {
		t.Errorf("unexpected failure: %+v", cfg.Logger())
	}
	if cfg.DB() != mockDBConfig {
		t.Errorf("unexpected failure: %+v", cfg.DB())
	}
	if cfg.GRPCServer() != mockGRPCServerConfig {
		t.Errorf("unexpected failure: %+v", cfg.GRPCServer())
	}
	if cfg.HTTPServer() != mockHTTPServerConfig {
		t.Errorf("unexpected failure: %+v", cfg.HTTPServer())
	}
}

func createConfig(t *testing.T, cfgFuncs ...interface{}) (*Config, error) {
	t.Helper()

	ctor := ConfigConstructors{
		NewLoggerConfig:     func(opts ...config.Option) (log.Config, error) { return nil, nil },
		NewDBConfig:         func(opts ...config.Option) (sql.Config, error) { return nil, nil },
		NewGRPCServerConfig: func(opts ...config.Option) (grpc.Config, error) { return nil, nil },
		NewHTTPServerConfig: func(opts ...config.Option) (http.Config, error) { return nil, nil },
	}
	for _, cfgFunc := range cfgFuncs {
		switch f := cfgFunc.(type) {
		case func(opts ...config.Option) (log.Config, error):
			ctor.NewLoggerConfig = f
		case func(opts ...config.Option) (sql.Config, error):
			ctor.NewDBConfig = f
		case func(opts ...config.Option) (grpc.Config, error):
			ctor.NewGRPCServerConfig = f
		case func(opts ...config.Option) (http.Config, error):
			ctor.NewHTTPServerConfig = f
		default:
			t.Errorf("Unexpected cfgCtor")
		}
	}
	return NewConfig(ctor)
}

func TestNewServerConfigFailures(t *testing.T) {
	t.Parallel()

	t.Run("Logger failure", func(t *testing.T) {
		t.Parallel()
		expectedErr := "NewLogger failure"
		cfg, err := createConfig(t, func(opts ...config.Option) (log.Config, error) { return nil, errors.New(expectedErr) })
		if cfg != nil || err.Error() != expectedErr {
			t.Errorf("Error expected: %+v %+v", cfg, err)
		}
	})

	t.Run("DB failure", func(t *testing.T) {
		t.Parallel()
		expectedErr := "NewDB failure"
		cfg, err := createConfig(t, func(opts ...config.Option) (sql.Config, error) { return nil, errors.New(expectedErr) })
		if cfg != nil || err.Error() != expectedErr {
			t.Errorf("Error expected: %+v %+v", cfg, err)
		}
	})

	t.Run("DB failure", func(t *testing.T) {
		t.Parallel()
		expectedErr := "NewGRPCServer failure"
		cfg, err := createConfig(t, func(opts ...config.Option) (grpc.Config, error) { return nil, errors.New(expectedErr) })
		if cfg != nil || err.Error() != expectedErr {
			t.Errorf("Error expected: %+v %+v", cfg, err)
		}
	})

	t.Run("DB failure", func(t *testing.T) {
		t.Parallel()
		expectedErr := "NewHTTPServer failure"
		cfg, err := createConfig(t, func(opts ...config.Option) (http.Config, error) { return nil, errors.New(expectedErr) })
		if cfg != nil || err.Error() != expectedErr {
			t.Errorf("Error expected: %+v %+v", cfg, err)
		}
	})
}
