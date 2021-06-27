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
	Options struct {
		prefix string
	}

	// Option provides options when collecting configuration information.
	Option interface {
		Apply(*Options)
	}

	option struct {
		f func(*Options)
	}
)

func (o *Options) Prefix() string {
	return o.prefix
}

func newOption(f func(*Options)) *option {
	return &option{f: f}
}

func (o *option) Apply(opts *Options) {
	o.f(opts)
}

// WithPrefix adds a prefix to the name of the enviroment variables being referenced.
func WithPrefix(prefix string) Option {
	return newOption(func(opts *Options) {
		if prefix != "" {
			opts.prefix = prefix + "_" + opts.prefix
		}
	})
}
