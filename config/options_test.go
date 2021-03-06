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

func TestOptionsWithPrefix(t *testing.T) {
	t.Parallel()

	t.Run("Test WithPrefix - empty prefix", func(t *testing.T) {
		opts := &config.Options{}

		config.WithPrefix("").Apply(opts)
		if opts.Prefix != "" {
			t.Error("prefix should be empty")
		}
	})

	t.Run("Test WithPrefix - non-empty prefix", func(t *testing.T) {
		opts := &config.Options{}

		config.WithPrefix("prefix").Apply(opts)
		if opts.Prefix != "prefix_" {
			t.Errorf("incorrect prefix: %s", opts.Prefix)
		}
	})
}
