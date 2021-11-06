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

package build // import "arcadium.dev/core/build"

import (
	"fmt"
	"runtime"
)

// Information holds the build information.
type Information struct {
	Name, Version, Branch, Commit, Date, Go string
}

// Info populates the build information.
func Info(n, v, b, c, d string) Information {
	return Information{
		Name:    n,
		Version: v,
		Branch:  b,
		Commit:  c,
		Date:    d,
		Go:      runtime.Version(),
	}
}

// Fields provides an intuitive way to add the build information to a log entry.
func (i Information) Fields() []interface{} {
	return []interface{}{
		"name", i.Name,
		"version", i.Version,
		"branch", i.Branch,
		"commit", i.Commit,
		"date", i.Date,
		"go", i.Go,
	}
}

// String provides the build information as a string.
func (i Information) String() string {
	return fmt.Sprintf("%s %s (branch: %s, commit: %s, date: %s, go: %s)",
		i.Name, i.Version, i.Branch, i.Commit, i.Date, i.Go,
	)
}
