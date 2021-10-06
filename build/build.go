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
	"os"
	"path/filepath"
	"runtime"
)

// Information holds the build information.
type Information map[string]interface{}

// Info populates the build information.
func Info(v, b, s, d string) Information {
	return map[string]interface{}{
		"name":    filepath.Base(os.Args[0]),
		"version": v,
		"branch":  b,
		"shasum":  s,
		"date":    d,
		"go":      runtime.Version(),
	}
}

// Fields provides an intuitive way to add the build information to a log entry.
func (i Information) Fields() map[string]interface{} {
	return i
}

// Version provides the build information as a version string.
func (i Information) Version() string {
	return fmt.Sprintf("%s %s (branch: %s, shasum: %s, date: %s, go: %s)",
		i["name"].(string),
		i["version"].(string),
		i["branch"].(string),
		i["shasum"].(string),
		i["date"].(string),
		i["go"].(string),
	)
}
