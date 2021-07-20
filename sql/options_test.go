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
	"database/sql"
	"testing"
)

func TestOptionsWithMigration(t *testing.T) {
	t.Parallel()

	t.Run("WithMigration with nil migration", func(t *testing.T) {
		t.Parallel()

		var opts options
		WithMigration(nil).apply(&opts)
		if opts.migration != nil {
			t.Errorf("Unexpected migration: %+v", opts.migration)
		}
	})

	t.Run("WithMigration with valid migration", func(t *testing.T) {
		t.Parallel()

		var opts options
		WithMigration(testMigration).apply(&opts)
		if opts.migration == nil {
			t.Error("Expected migration")
		}
	})
}

func testMigration(*sql.DB) error { return nil }
