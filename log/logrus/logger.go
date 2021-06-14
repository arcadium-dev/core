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
	"github.com/sirupsen/logrus"

	"arcadium.dev/core/log"
)

// Logger

type (
	Logger struct {
		entry *logrus.Entry
		level logrus.Level
	}
)

func New(cfg log.Config) *Logger {
	l := logrus.New()

	level, err := logrus.ParseLevel(cfg.Level())
	if err != nil {
		level = logrus.InfoLevel
	}
	l.SetLevel(level)

	return &Logger{
		entry: logrus.NewEntry(l),
		level: l.Level,
	}
}

func (l *Logger) WithField(key string, value interface{}) log.Logger {
	return &Logger{entry: l.entry.WithField(key, value), level: l.level}
}

func (l *Logger) WithFields(fields map[string]interface{}) log.Logger {
	return &Logger{entry: l.entry.WithFields(fields), level: l.level}
}

func (l *Logger) WithError(err error) log.Logger {
	return &Logger{entry: l.entry.WithError(err), level: l.level}
}

func (l *Logger) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *Logger) Debugln(args ...interface{}) {
	l.entry.Debugln(args...)
}

func (l *Logger) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *Logger) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *Logger) Infoln(args ...interface{}) {
	l.entry.Infoln(args...)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *Logger) Print(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *Logger) Println(args ...interface{}) {
	l.entry.Infoln(args...)
}

func (l *Logger) Printf(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *Logger) Warning(args ...interface{}) {
	l.entry.Warning(args...)
}

func (l *Logger) Warningln(args ...interface{}) {
	l.entry.Warningln(args...)
}

func (l *Logger) Warningf(format string, args ...interface{}) {
	l.entry.Warningf(format, args...)
}

func (l *Logger) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *Logger) Errorln(args ...interface{}) {
	l.entry.Errorln(args...)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *Logger) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *Logger) Fatalln(args ...interface{}) {
	l.entry.Fatalln(args...)
}

func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *Logger) V(level int) bool {
	return l.level <= logrus.Level(level)
}
