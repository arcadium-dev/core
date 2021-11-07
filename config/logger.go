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

package config // import "arcadium.dev/core/config/config

import (
	"github.com/kelseyhightower/envconfig"

	"arcadium.dev/core/errors"
)

const (
	logPrefix = "log"
)

type (
	// Logger provides the configuration settings for a logger.
	Logger interface {
		// Level returns the configured log level.
		Level() string
		// File returns the configured log file.
		File() string
		// Format returns the configured log format.
		Format() string
	}
)

// NewLogger returns the configuration of a logger.
func NewLogger(opts ...Option) (Logger, error) {
	o := &options{}
	for _, opt := range opts {
		opt.apply(o)
	}
	prefix := o.prefix + logPrefix

	config := struct {
		Level  string
		File   string
		Format string
	}{}
	if err := envconfig.Process(prefix, &config); err != nil {
		return nil, errors.Wrapf(err, "failed to load %s configuration", prefix)
	}
	return &logger{
		level:  config.Level,
		file:   config.File,
		format: config.Format,
	}, nil
}

type (
	logger struct {
		level  string // <PREFIX_>LOG_LEVEL
		file   string // <PREFIX_>LOG_FILE
		format string // <PREFIX_>LOG_FORMAT
	}
)

func (l *logger) Level() string {
	return l.level
}

func (l *logger) File() string {
	return l.file
}

func (l *logger) Format() string {
	return l.format
}
