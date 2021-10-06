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

package envconfig

import (
	"testing"

	"arcadium.dev/core/config"
)

func TestServer(t *testing.T) {
	t.Run("Empty Env", func(t *testing.T) {
		cfg := setupServer(t, config.Env(nil))

		if cfg.Addr() != "" || cfg.Cert() != "" || cfg.Key() != "" || cfg.CACert() != "" {
			t.Error("incorrect server config for an empty environment")
		}
	})

	t.Run("Full Env", func(t *testing.T) {
		cfg := setupServer(t, config.Env(map[string]string{
			"SERVER_ADDR":   "test_addr:42",
			"SERVER_CERT":   "/opt/cert.crt",
			"SERVER_KEY":    "/opt/key.crt",
			"SERVER_CACERT": "/opt/cacert.crt",
		}))

		if cfg.Addr() != "test_addr:42" || cfg.Cert() != "/opt/cert.crt" ||
			cfg.Key() != "/opt/key.crt" || cfg.CACert() != "/opt/cacert.crt" {
			t.Error("incorrect server config for a full environment")
		}
	})

	t.Run("Partial Env", func(t *testing.T) {
		cfg := setupServer(t, config.Env(map[string]string{
			"SERVER_CERT": "/opt/cert.crt",
			"SERVER_KEY":  "/opt/key.crt",
		}))

		if cfg.Addr() != "" || cfg.Cert() != "/opt/cert.crt" ||
			cfg.Key() != "/opt/key.crt" || cfg.CACert() != "" {
			t.Error("incorrect server config for a partial environment")
		}
	})

	t.Run("WithPrefix", func(t *testing.T) {
		cfg := setupServer(t, config.Env(map[string]string{
			"FANCY_SERVER_ADDR":   "test_addr:42",
			"FANCY_SERVER_CERT":   "/opt/cert.crt",
			"FANCY_SERVER_KEY":    "/opt/key.crt",
			"FANCY_SERVER_CACERT": "/opt/cacert.crt",
		}), config.WithPrefix("fancy"))

		if cfg.Addr() != "test_addr:42" || cfg.Cert() != "/opt/cert.crt" ||
			cfg.Key() != "/opt/key.crt" || cfg.CACert() != "/opt/cacert.crt" {
			t.Error("incorrect server config for a full environment")
		}
	})
}

func setupServer(t *testing.T, e config.Env, opts ...config.Option) *Server {
	t.Helper()

	e.Set(t)

	cfg, err := NewServer(opts...)
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	return cfg
}
