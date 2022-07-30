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

package config_test

import (
	"testing"

	"arcadium.dev/core/config"
)

func TestDB(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		_, err := config.NewDB()

		if err == nil {
			t.Errorf("expected an error")
		}
		expectedErr := "failed to load db configuration: required key DB_DRIVER missing value"
		if err.Error() != expectedErr {
			t.Errorf("\nExpected error: %s\nActual error:   %s", expectedErr, err)
		}
	})

	t.Run("success", func(t *testing.T) {
		expectedDSN := "postgresql://user@cockroach:26257/db?sslmode=verify-full&sslrootcert=%2Fetc%2Fcerts%2Fca.crt"
		t.Setenv("DB_DRIVER", "postgres")
		t.Setenv("DB_DSN", expectedDSN)
		cfg := setupDB(t)

		if cfg.Driver() != "postgres" {
			t.Errorf("Unexpected driver: %s", cfg.Driver())
		}
		if cfg.DSN() != expectedDSN {
			t.Errorf("\nExpected url: %s\nActual url:   %s", expectedDSN, cfg.DSN())
		}
	})
}

func setupDB(t *testing.T, opts ...config.Option) config.DB {
	t.Helper()

	cfg, err := config.NewDB(opts...)
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	return cfg
}
