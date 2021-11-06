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

func setup(t *testing.T) Information {
	t.Helper()
	info := Info("Testing", "Version", "Branch", "Commit", "Date")
	info.Go = "Go"
	return info
}

func TestFields(t *testing.T) {
	t.Parallel()

	fields := setup(t).Fields()

	for i, expected := range []string{
		"name", "Testing", "version", "Version", "branch", "Branch", "commit", "Commit", "date", "Date", "go", "Go",
	} {
		actual, ok := fields[1].(string)
		if !ok {
			t.Errorf("Failed type assertion")
		}
		if fields[i].(string) != expected {
			t.Errorf("Expected %s, Actual: %s", expected, actual)
		}
	}
}

func TestString(t *testing.T) {
	t.Parallel()

	actual := setup(t).String()

	expected := "Testing Version (branch: Branch, commit: Commit, date: Date, go: Go)"
	if actual != expected {
		t.Errorf("\nExpected: %s,\nActual:   %s", expected, actual)
	}
}
