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

package test // import "arcadium.dev/core/test"

import "testing"

// Env is a test helper to facilitate setting up environment variables.
type Env map[string]string

// Set places the environment variables in Env in the environment.  The
// environment variables are automatically cleaned up after test completion.
func (e Env) Set(t *testing.T) {
	for k, v := range e {
		t.Setenv(k, v)
	}
}
