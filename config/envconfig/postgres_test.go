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

func TestPostgres(t *testing.T) {
	t.Run("Empty Env", func(t *testing.T) {
		cfg := setupPostgres(t, config.Env(nil))

		expectedDSN := "pgx://host:port"
		if cfg.DSN() != expectedDSN {
			t.Errorf("\nExpected dsn: %s\nActual dsn:   %s", expectedDSN, cfg.DSN())
		}
	})

	t.Run("Full Env", func(t *testing.T) {
		cfg := setupPostgres(t, config.Env(map[string]string{
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
		cfg := setupPostgres(t, config.Env(map[string]string{
			"POSTGRES_DB":       "players",
			"POSTGRES_USER":     "arcadium",
			"POSTGRES_PASSWORD": "password",
			"POSTGRES_HOST":     "postgres",
			"POSTGRES_PORT":     "5432",
			"POSTGRES_SSLMODE":  "disable",
		}))

		expectedDSN := "postgres://arcadium:password@postgres:5432/players?sslmode=disable"
		if cfg.DSN() != expectedDSN {
			t.Errorf("\nExpected dsn: %s\nActual dsn    %s", expectedDSN, cfg.DSN())
		}
	})
}

func setupPostgres(t *testing.T, e config.Env) *Postgres {
	e.Set()
	defer e.Unset()

	cfg, err := NewPostgres()
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	return cfg
}
