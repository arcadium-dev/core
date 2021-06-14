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

package build

import (
	"testing"
)

func setup() Information {
	info := Info("version", "branch", "shasum", "date")
	info["name"] = "testing"
	info["go"] = "go"
	return info
}

func TestFields(t *testing.T) {
	info := setup()
	fields := info.Fields()

	version, ok := fields["version"].(string)
	if !ok || version != "version" {
		t.Error("version incorrect")
	}

	branch, ok := fields["branch"].(string)
	if !ok || branch != "branch" {
		t.Error("branch incorrect")
	}

	shasum, ok := fields["shasum"].(string)
	if !ok || shasum != "shasum" {
		t.Error("shasum incorrect")
	}

	date, ok := fields["date"].(string)
	if !ok || date != "date" {
		t.Error("date incorrect")
	}

	name, ok := fields["name"].(string)
	if !ok || name != "testing" {
		t.Error("date incorrect")
	}

	g, ok := fields["go"].(string)
	if !ok || g != "go" {
		t.Error("go incorrect")
	}
}

func TestVersion(t *testing.T) {
	info := setup()

	v := info.Version()

	if v != "testing version (branch: branch, shasum: shasum, date: date, go: go)" {
		t.Error("version incorrect")
	}
}
