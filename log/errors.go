// Copyright 2021-2022 arcadium.dev <info@arcadium.dev>
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

package log // import "arcadium.dev/core/log

import "errors"

var (
	// ErrInvalidLevel will be returned when the level given to the WirhLevel
	// option is invalid.
	ErrInvalidLevel = errors.New("invalid level")

	// ErrInvalidFormat will be returned when the format given to the WithFormat
	// options is invalid.
	ErrInvalidFormat = errors.New("invalid format")

	// ErrInvalidOutput will be returned when the output writer given to WithOuput
	// is nil.
	ErrInvalidOutput = errors.New("invalid output")
)
