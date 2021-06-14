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

func postgresSetup(t *testing.T, e config.Env) *Postgres {
	e.Set()
	defer e.Unset()

	cfg, err := NewPostgres()
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	return cfg
}

func TestPostgresFullEnv(t *testing.T) {
	expectedDSN := "dbname='db' user='user' password='password' host='host' port='port' connect_timeout='connect_timeout' sslmode='sslmode' sslcert='sslcert' sslkey='sslkey' sslrootcert='sslrootcert'"

	cfg := postgresSetup(t, config.Env(map[string]string{
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

	if cfg.DSN() != expectedDSN {
		t.Errorf("incorrect postgres DSN, expected %s, actual %s", expectedDSN, cfg.DSN())
	}
}

func TestPostgresPartialEnv(t *testing.T) {
	expectedDSN := "dbname='players' user='arcadium' password='password' host='postgres' port='5432' sslmode='disable'"

	cfg := postgresSetup(t, config.Env(map[string]string{
		"POSTGRES_DB":       "players",
		"POSTGRES_USER":     "arcadium",
		"POSTGRES_PASSWORD": "password",
		"POSTGRES_HOST":     "postgres",
		"POSTGRES_PORT":     "5432",
		"POSTGRES_SSLMODE":  "disable",
	}))

	if cfg.DSN() != expectedDSN {
		t.Errorf("incorrect postgres DSN, expected %s, actual %s", expectedDSN, cfg.DSN())
	}
}
