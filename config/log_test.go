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

func TestLog(t *testing.T) {
	t.Run("Empty Env", func(t *testing.T) {
		cfg := setupLogger(t)

		if cfg.Level() != "" || cfg.Format() != "" {
			t.Error("incorrect logging config for an empty environment")
		}
	})

	t.Run("Full Env", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "LEVEL")
		t.Setenv("LOG_FORMAT", "FORMAT")
		cfg := setupLogger(t)

		if cfg.Level() != "level" || cfg.Format() != "format" {
			t.Error("incorrect logging config for a full environment")
		}
	})

	t.Run("Partial Env", func(t *testing.T) {
		t.Setenv("LOG_LEVEL", "level")
		cfg := setupLogger(t)

		if cfg.Level() != "level" || cfg.Format() != "" {
			t.Error("incorrect logging config for a partial environment")
		}
	})

	t.Run("WithPrefix", func(t *testing.T) {
		t.Setenv("PREFIX_LOG_LEVEL", "level")
		t.Setenv("PREFIX_LOG_FORMAT", "format")
		cfg := setupLogger(t, config.WithPrefix("prefix"))

		if cfg.Level() != "level" || cfg.Format() != "format" {
			t.Error("incorrect logging config for a full environment")
		}
	})
}

func setupLogger(t *testing.T, opts ...config.Option) config.Logger {
	t.Helper()

	cfg, err := config.NewLogger(opts...)
	if err != nil {
		t.Fatalf("error occurred: %s", err)
	}
	return cfg
}
