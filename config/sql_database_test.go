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

func TestNewConfig(t *testing.T) {
	t.Run("Test default driver", func(t *testing.T) {
		cfg := setupSQLDatabase(t, test.Env(map[string]string{
			"POSTGRES_DB":   "db",
			"POSTGRES_HOST": "host",
		}))
		if cfg.DriverName() != "pgx" {
			t.Error("Incorrect sql database config for an empty environment")
		}
	})

	t.Run("Test postgres driver", func(t *testing.T) {
		cfg := setupSQLDatabase(t, test.Env(map[string]string{
			"SQL_DATABASE_DRIVER": "postgres",
			"POSTGRES_DB":         "db",
			"POSTGRES_HOST":       "host",
		}))
		if cfg.DriverName() != "postgres" {
			t.Error("Incorrect sql database config for a valid environment")
		}
	})

	t.Run("Test WithPrefix", func(t *testing.T) {
		cfg := setupSQLDatabase(t, test.Env(map[string]string{
			"FOO_SQL_DATABASE_DRIVER": "pgx",
			"FOO_POSTGRES_DB":         "db",
			"FOO_POSTGRES_HOST":       "host",
		}), WithPrefix("foo"))
		if cfg.DriverName() != "pgx" {
			t.Errorf("Incorrect sql database config for a valid environment: %s", cfg.DriverName())
		}
	})

	t.Run("Test unsupported driver", func(t *testing.T) {
		e := test.Env(map[string]string{
			"SQL_DATABASE_DRIVER": "mysql",
		})
		e.Set(t)

		cfg, err := NewSQLDatabase()
		if cfg != nil {
			t.Errorf("Unexpected sql database config")
		}
		if err == nil {
			t.Errorf("Error expected")
		}
		if err.Error() != "unsupported database driver: mysql" {
			t.Errorf("Unexpected error: %s", err)
		}
	})
}

func setupSQLDatabase(t *testing.T, e test.Env, opts ...Option) *SQLDatabase {
	t.Helper()

	e.Set(t)

	cfg, err := NewSQLDatabase(opts...)
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	return cfg
}
