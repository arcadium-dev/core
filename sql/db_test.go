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
	"testing"
)

func TestOpenSuccess(t *testing.T) {
	// FIXME
}

func TestOpenFailures(t *testing.T) {
	t.Parallel()

	t.Run("Test driver name failure", func(t *testing.T) {
		// FIXME
	})

	t.Run("Test dsn failure", func(t *testing.T) {
		// FIXME
	})

	t.Run("Test migration failure", func(t *testing.T) {
		// FIXME
	})
}

func TestConnectSuccess(t *testing.T) {
	// FIXME
}

func TestConnectFailure(t *testing.T) {
	// FIXME
}
