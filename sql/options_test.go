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

package sql

import (
	"io/fs"
	"testing"
)

func TestOptionsWithMigrator(t *testing.T) {
	t.Parallel()

	t.Run("WithMigrator with nil migrator", func(t *testing.T) {
		t.Parallel()

		var opts options
		WithMigrator(nil).apply(&opts)
		if opts.migrator != nil {
			t.Errorf("Unexpected migrator: %+v", opts.migrator)
		}
	})

	t.Run("WithMigrator with valid migrator", func(t *testing.T) {
		t.Parallel()

		m := &mockMigrator{}

		var opts options
		WithMigrator(m).apply(&opts)
		if opts.migrator != m {
			t.Errorf("Unexpected migrator: %+v", opts.migrator)
		}
	})
}

type mockMigrator struct{}

func (m mockMigrator) Migrate(int, fs.FS, DB) error { return nil }
