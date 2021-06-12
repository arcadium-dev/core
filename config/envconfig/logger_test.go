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

func loggerSetup(e config.Env, opts ...Option) (*Logger, error) {
	e.Set()
	defer e.Unset()
	return NewLogger(opts...)
}

func TestLoggerEmptyEnv(t *testing.T) {
	cfg, err := loggerSetup(config.Env(nil))
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	if cfg.Level() != "" || cfg.File() != "" || cfg.Format() != "" {
		t.Error("incorrect logging config for an empty environment")
	}
}

func TestLoggerFullEnv(t *testing.T) {
	cfg, err := loggerSetup(config.Env(map[string]string{
		"LOG_LEVEL":  "level",
		"LOG_FILE":   "file",
		"LOG_FORMAT": "format",
	}))
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	if cfg.Level() != "level" || cfg.File() != "file" || cfg.Format() != "format" {
		t.Error("incorrect logging config for a full environment")
	}
}

func TestLoggerWithPrefux(t *testing.T) {
	cfg, err := loggerSetup(config.Env(map[string]string{
		"PREFIX_LOG_LEVEL":  "level",
		"PREFIX_LOG_FILE":   "file",
		"PREFIX_LOG_FORMAT": "format",
	}), WithPrefix("prefix"))
	if err != nil {
		t.Errorf("error occurred: %s", err)
	}
	if cfg.Level() != "level" || cfg.File() != "file" || cfg.Format() != "format" {
		t.Error("incorrect logging config for a full environment")
	}
}
