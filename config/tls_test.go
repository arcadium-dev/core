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

package config

import (
	"crypto/tls"
	"testing"
)

func TestNewTLS(t *testing.T) {
	t.Parallel()

	t.Run("Without options, bad cert", func(t *testing.T) {
		t.Parallel()

		var mockCfg = mockTLS{cert: "bad cert"}
		cfg, err := NewTLS(mockCfg)
		if cfg != nil {
			t.Errorf("Unexpected cfg: %+v", cfg)
		}
		if err == nil {
			t.Errorf("Expected an error")
		}
		if err.Error() != "open bad cert: no such file or directory: Failed to load server certificate" {
			t.Errorf("Unexpected err: %s", err)
		}
	})

	t.Run("Without options, bad key", func(t *testing.T) {
		t.Parallel()

		var mockCfg = mockTLS{key: "bad key"}
		cfg, err := NewTLS(mockCfg)
		if cfg != nil {
			t.Errorf("Unexpected cfg: %+v", cfg)
		}
		if err == nil {
			t.Errorf("Expected an error")
		}
		if err.Error() != "open bad key: no such file or directory: Failed to load server certificate" {
			t.Errorf("Unexpected err: %s", err)
		}
	})

	t.Run("Without options, success", func(t *testing.T) {
		t.Parallel()

		var mockCfg = mockTLS{}
		cfg, err := NewTLS(mockCfg)
		if cfg == nil {
			t.Errorf("Expected a cfg")
		}
		if err != nil {
			t.Errorf("Unexpected err: %s", err)
		}
	})

	t.Run("WithMTLS option, bad cacert", func(t *testing.T) {
		t.Parallel()

		var mockCfg = mockTLS{cacert: "bad cacert"}
		cfg, err := NewTLS(mockCfg, WithMTLS())
		if cfg != nil {
			t.Errorf("Unexpected cfg: %+v", cfg)
		}
		if err == nil {
			t.Errorf("Expected an error")
		}
		if err.Error() != "open bad cacert: no such file or directory: Failed to load the client CA certificate" {
			t.Errorf("Unexpected err: %s", err)
		}
	})

	t.Run("WithMTLS option, success", func(t *testing.T) {
		t.Parallel()

		var mockCfg = mockTLS{}
		cfg, err := NewTLS(mockCfg, WithMTLS())
		if cfg == nil {
			t.Errorf("Expected a cfg")
		}
		if err != nil {
			t.Errorf("Unexpected err: %s", err)
		}
	})
}

func TestWithMTLS(t *testing.T) {
	t.Parallel()

	cfg := &tls.Config{}
	WithMTLS().apply(cfg)

	if cfg.ClientAuth != tls.RequireAndVerifyClientCert {
		t.Errorf("Unexpected ClientAuth: %+v", cfg.ClientAuth)
	}
}

type mockTLS struct {
	cert, key, cacert string
}

func (m mockTLS) Cert() string {
	if m.cert == "" {
		return "../insecure/cert.pem"
	}
	return m.cert
}

func (m mockTLS) Key() string {
	if m.key == "" {
		return "../insecure/key.pem"
	}
	return m.key
}

func (m mockTLS) CACert() string {
	if m.cacert == "" {
		return "../insecure/rootCA.pem"
	}
	return m.cacert
}
