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
)

func TestServer(t *testing.T) {
	t.Run("without prefix", func(t *testing.T) {
		t.Setenv("SERVER_ADDR", "test_addr:42")
		cfg := setupServer(t)

		if cfg.Addr() != "test_addr:42" {
			t.Error("incorrect server config")
		}
	})

	t.Run("with prefix", func(t *testing.T) {
		t.Setenv("FANCY_SERVER_ADDR", "test_addr:42")
		cfg := setupServer(t, WithPrefix("fancy"))

		if cfg.Addr() != "test_addr:42" {
			t.Error("incorrect server config")
		}
	})
}

func setupServer(t *testing.T, opts ...Option) Server {
	t.Helper()

	cfg, err := NewServer(opts...)
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	return cfg
}
