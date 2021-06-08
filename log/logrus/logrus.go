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
	"os"

	"github.com/pkg/errors"

	"arcadium.dev/core/log"
)

type (
	Logger struct {
		*logger
		file *os.File
	}
)

func New(cfg log.Config) (log.Logger, error) {
	var (
		err     error
		file    *os.File
		options []Option
	)

	if filename := cfg.File(); filename != "" {
		file, err = os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to open log file: %s", filename)
		}
		options = append(options, WithOutput(file))
	}

	if level := cfg.Level(); level != "" {
		options = append(options, WithLevel(level))
	}

	if format := cfg.Format(); format != "" {
		options = append(options, WithFormat(format))
	}

	return &Logger{
		logger: newLogger(options...),
		file:   file,
	}, nil
}

func (l *Logger) Close() error {
	if l.file == nil {
		return nil
	}
	return l.file.Close()
}
