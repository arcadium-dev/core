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

package errors

import (
	"errors"
	"fmt"
)

var (
	// New is an alias for errors.New
	New = errors.New

	// Errorf is an alias for fmt.Errorf
	Errorf = fmt.Errorf

	// Unwrap is an alias for errors.Unwrap
	Unwrap = errors.Unwrap

	// Is is an alias for errors.Is
	Is = errors.Is

	// As is an alias for errors.As
	As = errors.As
)

// Wrap returns an error wrapped with the supplied message.
// If the given error is nil, Wrap returns nil.
func Wrap(err error, msg string) error {
	if err == nil {
		return nil
	}
	return Errorf("%w: %s", err, msg)
}

// Wrap returns an error wrapped with the format specifier.
// If the given error is nil, Wrap returns nil.
func Wrapf(e error, format string, a ...interface{}) error {
	return Wrap(e, fmt.Sprintf(format, a...))
}
