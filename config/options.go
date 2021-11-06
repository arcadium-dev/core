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

package config // import "arcadium.dev/core/config

type (
	// Option provides options when loading configuration information.
	Option interface {
		apply(*options)
	}
)

// WithPrefix adds a prefix to the name of the enviroment variables being referenced.
func WithPrefix(prefix string) Option {
	return newOption(func(opts *options) {
		if prefix != "" {
			opts.prefix = prefix + "_" + opts.prefix
		}
	})
}

type (
	options struct {
		prefix string
	}

	option struct {
		f func(*options)
	}
)

func (o *option) apply(opts *options) {
	o.f(opts)
}

func newOption(f func(*options)) *option {
	return &option{f: f}
}
