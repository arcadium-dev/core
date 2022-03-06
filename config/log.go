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

package config // import "arcadium.dev/core/config/config

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type (
	// Logger holds the configuration information for a logger.
	Logger struct {
		level  string
		format string
	}
)

const (
	logPrefix = "log"
)

// NewLogger returns the configuration of a logger.
func NewLogger(opts ...Option) (Logger, error) {
	o := &Options{}
	for _, opt := range opts {
		opt.Apply(o)
	}
	prefix := o.Prefix + logPrefix

	config := struct {
		Level  string
		Format string
	}{}
	if err := envconfig.Process(prefix, &config); err != nil {
		return Logger{}, fmt.Errorf("failed to load %s configuration: %w", prefix, err)
	}
	return Logger{
		level:  strings.TrimSpace(strings.ToLower(config.Level)),
		format: strings.TrimSpace(strings.ToLower(config.Format)),
	}, nil
}

// Level returns the logging level for the logger. The value is set from the
// <PREFIX_>LOG_LEVEL environment variable.
func (l Logger) Level() string {
	return l.level
}

// Format returns the logging format for the logger. The value is set from the
// <PREFIX_>LOG_FORMAT environment variable.
func (l Logger) Format() string {
	return l.format
}
