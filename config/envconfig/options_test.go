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
)

func TestOptions(t *testing.T) {
	opts := &options{}

	WithPrefix("").apply(opts)
	if opts.prefix != "" {
		t.Error("prefix should be empty")
	}

	WithPrefix("prefix").apply(opts)
	if opts.prefix != "prefix_" {
		t.Errorf("incorrect prefix: %s", opts.prefix)
	}
}
