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

package envconfig // import "arcadium.dev/core/config/envconfig

import (
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"
)

const (
	logPrefix = "log"
)

type (
	// Logger holds the configuration settings for a logger.
	Logger struct {
		level  string // <PREFIX_>LOG_LEVEL
		file   string // <PREFIX_>LOG_FILE
		format string // <PREFIX_>LOG_FORMAT
	}
)

// NewLogger returns the configuration of a logger.
func NewLogger(opts ...Option) (*Logger, error) {
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
	return &Logger{
		level:  config.Level,
		file:   config.File,
		format: config.Format,
	}, nil
}

// Level returns the configured log level.
func (l *Logger) Level() string {
	return l.level
}

// File returns the configured log file.
func (l *Logger) File() string {
	return l.file
}

// Format returns the configured log format.
func (l *Logger) Format() string {
	return l.format
}
