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

package logrus // import "arcadium.dev/core/log/logrus"

import (
	"io"

	"github.com/sirupsen/logrus"
)

type (
	options struct {
		level     *logrus.Level
		output    io.Writer
		formatter logrus.Formatter
	}

	Option interface {
		apply(*options)
	}

	option struct {
		f func(*options)
	}
)

func newOption(f func(opt *options)) *option {
	return &option{f: f}
}

func (o *option) apply(opt *options) {
	o.f(opt)
}

func WithLevel(level string) Option {
	return newOption(func(opts *options) {
		l, err := logrus.ParseLevel(level)
		if err != nil {
			return
		}
		opts.level = new(logrus.Level)
		*opts.level = l
	})
}

func WithOutput(out io.Writer) Option {
	return newOption(func(opts *options) {
		opts.output = out
	})
}

func WithFormat(f string) Option {
	return newOption(func(opts *options) {
		switch f {
		case "json":
			opts.formatter = &logrus.JSONFormatter{}
		default:
			opts.formatter = &logrus.TextFormatter{}
		}
	})
}
