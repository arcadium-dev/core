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

func TestLogger(t *testing.T) {
	t.Run("Empty Env", func(t *testing.T) {
		cfg := setupLogger(t, config.Env(nil))

		if cfg.Level() != "" || cfg.File() != "" || cfg.Format() != "" {
			t.Error("incorrect logging config for an empty environment")
		}
	})

	t.Run("Full Env", func(t *testing.T) {
		cfg := setupLogger(t, config.Env(map[string]string{
			"LOG_LEVEL":  "level",
			"LOG_FILE":   "file",
			"LOG_FORMAT": "format",
		}))

		if cfg.Level() != "level" || cfg.File() != "file" || cfg.Format() != "format" {
			t.Error("incorrect logging config for a full environment")
		}
	})

	t.Run("Partial Env", func(t *testing.T) {
		cfg := setupLogger(t, config.Env(map[string]string{
			"LOG_LEVEL": "level",
		}))

		if cfg.Level() != "level" || cfg.File() != "" || cfg.Format() != "" {
			t.Error("incorrect logging config for a partial environment")
		}
	})

	t.Run("WithPrefix", func(t *testing.T) {
		cfg := setupLogger(t, config.Env(map[string]string{
			"PREFIX_LOG_LEVEL":  "level",
			"PREFIX_LOG_FILE":   "file",
			"PREFIX_LOG_FORMAT": "format",
		}), config.WithPrefix("prefix"))

		if cfg.Level() != "level" || cfg.File() != "file" || cfg.Format() != "format" {
			t.Error("incorrect logging config for a full environment")
		}
	})
}

func setupLogger(t *testing.T, e config.Env, opts ...config.Option) *Logger {
	t.Helper()
	e.Set(t)

	cfg, err := NewLogger(opts...)
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	return cfg
}
