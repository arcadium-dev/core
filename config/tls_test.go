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

	"arcadium.dev/core/test"
)

const (
	goodCert   = "../test/insecure/cert.pem"
	goodKey    = "../test/insecure/key.pem"
	goodCACert = "../test/insecure/rootCA.pem"

	badCert   = "bad cert"
	badKey    = "bad key"
	badCACert = "bad cacert"
)

func TestTLS(t *testing.T) {
	t.Run("Minimal Env", func(t *testing.T) {
		cfg := setupTLS(t, test.Env(map[string]string{}))

		if cfg.Cert() != "" || cfg.Key() != "" || cfg.CACert() != "" {
			t.Error("incorrect tls config for an empty environment")
		}
	})

	t.Run("Full Env", func(t *testing.T) {
		cfg := setupTLS(t, test.Env(map[string]string{
			"TLS_CERT":   "/opt/cert.crt",
			"TLS_KEY":    "/opt/key.crt",
			"TLS_CACERT": "/opt/cacert.crt",
		}))

		if cfg.Cert() != "/opt/cert.crt" || cfg.Key() != "/opt/key.crt" || cfg.CACert() != "/opt/cacert.crt" {
			t.Error("incorrect tls config for a full environment")
		}
	})

	t.Run("WithPrefix", func(t *testing.T) {
		cfg := setupTLS(t, test.Env(map[string]string{
			"FANCY_TLS_CERT":   "/opt/cert.crt",
			"FANCY_TLS_KEY":    "/opt/key.crt",
			"FANCY_TLS_CACERT": "/opt/cacert.crt",
		}), WithPrefix("fancy"))

		if cfg.Cert() != "/opt/cert.crt" || cfg.Key() != "/opt/key.crt" || cfg.CACert() != "/opt/cacert.crt" {
			t.Error("incorrect tls config for a full environment")
		}
	})
}

func TestTLSConfig(t *testing.T) {
	t.Run("Without options, bad cert", func(t *testing.T) {
		var cfg = TLS{
			cert: badCert,
			key:  goodKey,
		}
		tlsCfg, err := cfg.TLSConfig()
		if tlsCfg != nil {
			t.Errorf("Unexpected cfg: %+v", cfg)
		}
		if err == nil {
			t.Errorf("Expected an error")
		}
		expected := "failed to load tls certificate: open bad cert: no such file or directory"
		if err.Error() != expected {
			t.Errorf("\nExpected error: %s\nActual error:   %s", expected, err)
		}
	})

	t.Run("Without options, bad key", func(t *testing.T) {
		var cfg = TLS{
			cert: goodCert,
			key:  badKey,
		}
		tlsCfg, err := cfg.TLSConfig()
		if tlsCfg != nil {
			t.Errorf("Unexpected cfg: %+v", cfg)
		}
		if err == nil {
			t.Errorf("Expected an error")
		}
		expected := "failed to load tls certificate: open bad key: no such file or directory"
		if err.Error() != expected {
			t.Errorf("\nExpected error: %s\nActual error:   %s", expected, err)
		}
	})

	t.Run("Without options, success", func(t *testing.T) {
		var cfg = TLS{
			cert: goodCert,
			key:  goodKey,
		}
		tlsCfg, err := cfg.TLSConfig()
		if tlsCfg == nil {
			t.Errorf("Expected a tls cfg")
		}
		if err != nil {
			t.Errorf("Unexpected err: %s", err)
		}
	})

	t.Run("WithMTLS option, bad cacert", func(t *testing.T) {
		var cfg = TLS{
			cert:   goodCert,
			key:    goodKey,
			cacert: badCACert,
		}
		tlsCfg, err := cfg.TLSConfig(WithMTLS())
		if tlsCfg != nil {
			t.Errorf("Unexpected cfg: %+v", cfg)
		}
		if err == nil {
			t.Errorf("Expected an error")
		}
		expected := "failed to load the client CA certificate: open bad cacert: no such file or directory"
		if err.Error() != expected {
			t.Errorf("\nExpected error: %s\nActual error:   %s", expected, err)
		}
	})

	t.Run("WithMTLS option, no cacert, success (assumes ca cert available from system)", func(t *testing.T) {
		var cfg = TLS{
			cert: goodCert,
			key:  goodKey,
		}
		tlsCfg, err := cfg.TLSConfig(WithMTLS())
		if tlsCfg == nil {
			t.Errorf("Expected a cfg")
		}
		if err != nil {
			t.Errorf("Unexpected err: %s", err)
		}
	})

	t.Run("WithMTLS option, cacert available, success", func(t *testing.T) {
		var cfg = TLS{
			cert:   goodCert,
			key:    goodKey,
			cacert: goodCACert,
		}
		tlsCfg, err := cfg.TLSConfig(WithMTLS())
		if tlsCfg == nil {
			t.Errorf("Expected a cfg")
		}
		if err != nil {
			t.Errorf("Unexpected err: %s", err)
		}
	})
}

func TestWithMTLS(t *testing.T) {
	cfg := &tls.Config{}
	WithMTLS().apply(cfg)

	if cfg.ClientAuth != tls.RequireAndVerifyClientCert {
		t.Errorf("Unexpected ClientAuth: %+v", cfg.ClientAuth)
	}
}

func setupTLS(t *testing.T, e test.Env, opts ...Option) TLS {
	t.Helper()
	e.Set(t)

	cfg, err := NewTLS(opts...)
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	return cfg
}
