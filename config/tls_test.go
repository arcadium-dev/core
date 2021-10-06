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

const (
	goodCert   = "../test/insecure/cert.pem"
	goodKey    = "../test/insecure/key.pem"
	goodCACert = "../test/insecure/rootCA.pem"

	badCert   = "bad cert"
	badKey    = "bad key"
	badCACert = "bad cacert"
)

func TestNewTLS(t *testing.T) {
	t.Parallel()

	t.Run("Without options, bad cert", func(t *testing.T) {
		t.Parallel()

		var mockCfg = mockTLS{
			cert: badCert,
			key:  goodKey,
		}
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

		var mockCfg = mockTLS{
			cert: goodCert,
			key:  badKey,
		}
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

		var mockCfg = mockTLS{
			cert: goodCert,
			key:  goodKey,
		}
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

		var mockCfg = mockTLS{
			cert:   goodCert,
			key:    goodKey,
			cacert: badCACert,
		}
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

	t.Run("WithMTLS option, no cacert, success (assumes ca cert available from system)", func(t *testing.T) {
		t.Parallel()

		var mockCfg = mockTLS{
			cert: goodCert,
			key:  goodKey,
		}
		cfg, err := NewTLS(mockCfg, WithMTLS())
		if cfg == nil {
			t.Errorf("Expected a cfg")
		}
		if err != nil {
			t.Errorf("Unexpected err: %s", err)
		}
	})

	t.Run("WithMTLS option, cacert available, success", func(t *testing.T) {
		t.Parallel()

		var mockCfg = mockTLS{
			cert:   goodCert,
			key:    goodKey,
			cacert: goodCACert,
		}
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

type (
	mockTLS struct {
		cert, key, cacert string
	}
)

func (m mockTLS) Cert() string   { return m.cert }
func (m mockTLS) Key() string    { return m.key }
func (m mockTLS) CACert() string { return m.cacert }
