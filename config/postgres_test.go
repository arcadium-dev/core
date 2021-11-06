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
	"testing"

	"arcadium.dev/core/test"
)

func TestPostgres(t *testing.T) {
	t.Run("Minimal Env", func(t *testing.T) {
		cfg := setupPostgres(t, test.Env(map[string]string{
			"POSTGRES_DB":   "db",
			"POSTGRES_HOST": "host",
		}))

		expectedDSN := "postgres://host/db?sslmode=verify-full"
		if cfg.DSN() != expectedDSN {
			t.Errorf("\nExpected dsn: %s\nActual dsn:   %s", expectedDSN, cfg.DSN())
		}
	})

	t.Run("Full Env", func(t *testing.T) {
		cfg := setupPostgres(t, test.Env(map[string]string{
			"POSTGRES_DB":              "db",
			"POSTGRES_USER":            "user",
			"POSTGRES_PASSWORD":        "password",
			"POSTGRES_HOST":            "host",
			"POSTGRES_PORT":            "port",
			"POSTGRES_CONNECT_TIMEOUT": "connect_timeout",
			"POSTGRES_SSLMODE":         "sslmode",
			"POSTGRES_SSLCERT":         "sslcert",
			"POSTGRES_SSLKEY":          "sslkey",
			"POSTGRES_SSLROOTCERT":     "sslrootcert",
		}))

		expectedDSN := "postgres://user:password@host:port/db?connect_timeout=connect_timeout&sslcert=sslcert&sslkey=sslkey&sslmode=sslmode&sslrootcert=sslrootcert"
		if cfg.DSN() != expectedDSN {
			t.Errorf("\nExpected dsn: %s\nActual dsn:   %s", expectedDSN, cfg.DSN())
		}
	})

	t.Run("Partial Env", func(t *testing.T) {
		cfg := setupPostgres(t, test.Env(map[string]string{
			"FOO_POSTGRES_DB":       "players",
			"FOO_POSTGRES_USER":     "arcadium",
			"FOO_POSTGRES_PASSWORD": "password",
			"FOO_POSTGRES_HOST":     "postgres",
			"FOO_POSTGRES_PORT":     "5432",
			"FOO_POSTGRES_SSLMODE":  "disable",
		}), WithPrefix("foo"))

		expectedDSN := "postgres://arcadium:password@postgres:5432/players?sslmode=disable"
		if cfg.DSN() != expectedDSN {
			t.Errorf("\nExpected dsn: %s\nActual dsn    %s", expectedDSN, cfg.DSN())
		}
	})
}

func TestPostgresFailure(t *testing.T) {
	t.Run("Empty Env", func(t *testing.T) {
		cfg, err := NewPostgres()

		if cfg != nil {
			t.Errorf("expected a nil cfg: %+v", cfg)
		}
		if err == nil {
			t.Errorf("expected an error")
		}
		expectedErr := "required key POSTGRES_DB missing value: failed to load postgres configuration"
		if err.Error() != expectedErr {
			t.Errorf("\nExpected error: %s\nActual error  %s", expectedErr, err)
		}
	})

	t.Run("Missing Host", func(t *testing.T) {
		e := test.Env(map[string]string{
			"POSTGRES_DB": "players",
		})
		e.Set(t)

		cfg, err := NewPostgres()

		if cfg != nil {
			t.Errorf("expected a nil cfg: %+v", cfg)
		}
		if err == nil {
			t.Errorf("expected an error")
		}
		expectedErr := "required key POSTGRES_HOST missing value: failed to load postgres configuration"
		if err.Error() != expectedErr {
			t.Errorf("\nExpected error: %s\nActual error  %s", expectedErr, err)
		}
	})
}

func setupPostgres(t *testing.T, e test.Env, opts ...Option) *Postgres {
	t.Helper()

	e.Set(t)

	cfg, err := NewPostgres(opts...)
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	return cfg
}
