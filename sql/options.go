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

package sql // import "arcadium.dev/core/sql"

type (
	options struct {
		logger Logger
	}

	// Option sets options such as instumenting the DB with logging.
	Option interface {
		apply(*options)
	}

	// option wraps a function that modifies options into an implementation
	// of the Option interface.
	option struct {
		f func(*options)
	}
)

func newOption(f func(*options)) *option {
	return &option{f: f}
}

func (o *option) apply(opts *options) {
	o.f(opts)
}

// WithLogger returns an Option that instruments the DB with logging.
func WithLogger(logger Logger) Option {
	return newOption(func(opts *options) { opts.logger = logger })
}
