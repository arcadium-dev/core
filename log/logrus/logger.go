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

package logrus // import "arcadium.dev/core/log/logrus"

import (
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"arcadium.dev/core/log"
)

// Logger

type (
	Logger struct {
		*logger
		file *os.File
	}
)

func New(cfg log.Config) (*Logger, error) {
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

// logger

type (
	logger struct {
		entry *logrus.Entry
		level logrus.Level
	}
)

func newLogger(opts ...Option) *logger {
	l := logrus.New()

	o := &options{}
	for _, opt := range opts {
		opt.apply(o)
	}
	if o.level != nil {
		l.SetLevel(*o.level)
	}
	if o.output != nil {
		l.SetOutput(o.output)
	}
	if o.formatter != nil {
		l.SetFormatter(o.formatter)
	}

	return &logger{
		entry: logrus.NewEntry(l),
		level: l.Level,
	}
}

func (l *logger) WithField(key string, value interface{}) log.Logger {
	return &logger{entry: l.entry.WithField(key, value), level: l.level}
}

func (l *logger) WithFields(fields map[string]interface{}) log.Logger {
	return &logger{entry: l.entry.WithFields(fields), level: l.level}
}

func (l *logger) WithError(err error) log.Logger {
	return &logger{entry: l.entry.WithError(err), level: l.level}
}

func (l *logger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *logger) Debugln(args ...interface{}) {
	l.entry.Debugln(args...)
}

func (l *logger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *logger) Infoln(args ...interface{}) {
	l.entry.Infoln(args...)
}

func (l *logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logger) Print(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *logger) Println(args ...interface{}) {
	l.entry.Infoln(args...)
}

func (l *logger) Printf(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logger) Warning(args ...interface{}) {
	l.entry.Warning(args...)
}

func (l *logger) Warningln(args ...interface{}) {
	l.entry.Warningln(args...)
}

func (l *logger) Warningf(format string, args ...interface{}) {
	l.entry.Warningf(format, args...)
}

func (l *logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *logger) Errorln(args ...interface{}) {
	l.entry.Errorln(args...)
}

func (l *logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *logger) Fatalln(args ...interface{}) {
	l.entry.Fatalln(args...)
}

func (l *logger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logger) V(level int) bool {
	return l.level <= logrus.Level(level)
}
