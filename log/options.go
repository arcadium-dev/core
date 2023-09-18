// Copyright 2023 arcadium.dev <info@arcadium.dev>
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

package log // import "arcadium.dev/core/log"

import (
	"io"

	"github.com/rs/zerolog"
)

type (
	// Option provides for Logger configuration.
	Option interface {
		apply(*options)
	}
)

// AsDefault sets the DefaultLogger.
func AsDefault() Option {
	return newOption(func(opts *options) {
		opts.asDefault = true
	})
}

// WithLevel allows the level to be configured. The default level is LevelInfo.
func WithLevel(level zerolog.Level) Option {
	return newOption(func(opts *options) {
		opts.level = level
	})
}

// WithLevelFieldName allows the customization of the level name. It defaults
// to "level".
func WithLevelFieldName(name string) Option {
	return newOption(func(opts *options) {
		opts.levelFieldName = name
	})
}

// WithMessageFieldName allows the customization of the message name. It
// defaults to "message".
func WithMessageFieldName(name string) Option {
	return newOption(func(opts *options) {
		opts.messageFieldName = name
	})
}

// WithOutput allows the format to be configured. The default writer is
// os.Stdout.
func WithOutput(writer io.Writer) Option {
	return newOption(func(opts *options) {
		opts.writer = writer
	})
}

// WithoutTimestamp disables the use of a timestamp for logs.  Useful for unit
// tests.
func WithoutTimestamp() Option {
	return newOption(func(opts *options) {
		opts.timestamped = false
	})
}

// WithTimestampFieldName allows the customization of the timestamp name. It
// defaults to "time".
func WithTimestampFieldName(name string) Option {
	return newOption(func(opts *options) {
		opts.timestampFieldName = name
	})
}

type (
	options struct {
		asDefault          bool
		level              zerolog.Level
		levelFieldName     string
		messageFieldName   string
		timestamped        bool
		timestampFieldName string
		writer             io.Writer
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
