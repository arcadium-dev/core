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

import "io"

type (
	// Option provides for Logger configuration.
	Option interface {
		apply(*options)
	}
)

// WithLevel allows the level to be configured. The default level is LevelInfo.
func WithLevel(level Level) Option {
	return newOption(func(opts *options) {
		opts.level = level
	})
}

// WithFormat allows the format to be configured. The default format is
// FormatLogfmt.
func WithFormat(format Format) Option {
	return newOption(func(opts *options) {
		opts.format = format
	})
}

// WithOutput allows the format to be configured. The default writer is
// os.Stdout.
func WithOutput(writer io.Writer) Option {
	return newOption(func(opts *options) {
		opts.writer = writer
	})
}

// WithoutTimestamp disables the use of a timestamp for logs.
// Useful for unit tests.
func WithoutTimestamp() Option {
	return newOption(func(opts *options) {
		opts.timestamped = false
	})
}

// As default sets the DefaultLogger.
func AsDefault() Option {
	return newOption(func(opts *options) {
		opts.asDefault = true
	})
}

type (
	options struct {
		level       Level
		format      Format
		writer      io.Writer
		timestamped bool
		asDefault   bool
	}

	option struct {
		f func(*options)
	}

	contextKey int
)

const (
	loggerContextKey = contextKey(iota + 1)
)

var (
	timestamped = true // Setting to false disables timestamps for unit testing.
)

func newOption(f func(*options)) *option {
	return &option{f: f}
}

func (o *option) apply(opts *options) {
	o.f(opts)
}
