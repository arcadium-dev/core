// Copyright 2021 Ian Cahoon <icahoon@gmail.com>
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

func setupSQLDatabase(t *testing.T, e config.Env, opts ...Option) *SQLDatabase {
	e.Set()
	defer e.Unset()

	cfg, err := NewSQLDatabase(opts...)
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	return cfg
}

func TestSQLDatabaseDefaultDriver(t *testing.T) {
	cfg := setupSQLDatabase(t, config.Env(nil))

	if cfg.DriverName() != "postgres" {
		t.Error("incorrect sql database config for an empty environment")
	}
}

func TestSQLDatabaseValidDriver(t *testing.T) {
	cfg := setupSQLDatabase(t, config.Env(map[string]string{
		"SQL_DATABASE_DRIVER": "postgres",
	}))

	if cfg.DriverName() != "postgres" {
		t.Error("incorrect sql database config for a valid environment")
	}
}

func TestSQLDatabaseWithPrefix(t *testing.T) {
	cfg := setupSQLDatabase(t, config.Env(map[string]string{
		"PLAYERS_DATABASE_DRIVER": "postgres",
	}), WithPrefix("players"))

	if cfg.DriverName() != "postgres" {
		t.Error("incorrect sql database config for a valid environment")
	}
}

func TestSQLDatabaseUnsupportedDriver(t *testing.T) {
	e := config.Env(map[string]string{
		"SQL_DATABASE_DRIVER": "mysql",
	})
	e.Set()
	defer e.Unset()

	cfg, err := NewSQLDatabase()
	if cfg != nil {
		t.Errorf("unexpected sql database config")
	}
	if err == nil || err.Error() != "unsupported database driver: mysql" {
		t.Errorf("error expected")
	}
}
