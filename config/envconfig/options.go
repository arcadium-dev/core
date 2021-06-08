// Copyright 2021 Ian Cahoon <icahoon@gmail.com>
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

package envconfig // import "arcadium.dev/core/config/envconfig

type (
	options struct {
		prefix string
	}

	// Option ... FIXME
	Option interface {
		apply(*options)
	}

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

// WithPrefix ... FIXME
func WithPrefix(prefix string) Option {
	return newOption(func(opts *options) {
		if prefix != "" {
			opts.prefix = prefix + "_" + opts.prefix
		}
	})
}
