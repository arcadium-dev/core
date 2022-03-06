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

package config

import (
	"testing"

	"arcadium.dev/core/test"
)

func TestSQL(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		_, err := NewSQL()

		if err == nil {
			t.Errorf("expected an error")
		}
		expectedErr := "failed to load sql configuration: required key SQL_URL missing value"
		if err.Error() != expectedErr {
			t.Errorf("\nExpected error: %s\nActual error  %s", expectedErr, err)
		}
	})

	t.Run("defaults", func(t *testing.T) {
		cfg := setupSQL(t, test.Env(map[string]string{
			"SQL_URL": "postgresql://user@cockroach:16567/db",
		}))

		if cfg.Driver() != "pgx" {
			t.Errorf("Expected: %s, Actual: %s", "pgx", cfg.Driver())
		}
	})

	t.Run("success", func(t *testing.T) {
		expectedURL := "postgresql://user@cockroach:26257/db?sslmode=verify-full&sslrootcert=%2Fetc%2Fcerts%2Fca.crt"
		cfg := setupSQL(t, test.Env(map[string]string{
			"SQL_DRIVER": "postgres",
			"SQL_URL":    expectedURL,
		}))

		if cfg.Driver() != "postgres" {
			t.Errorf("Unexpected driver: %s", cfg.Driver())
		}
		if cfg.URL() != expectedURL {
			t.Errorf("\nExpected url: %s\nActual url:   %s", expectedURL, cfg.URL())
		}
	})
}

func setupSQL(t *testing.T, e test.Env, opts ...Option) SQL {
	t.Helper()
	e.Set(t)

	cfg, err := NewSQL(opts...)
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	return cfg
}
