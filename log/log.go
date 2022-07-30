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

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

var (
	DefaultLogger Logger
)

func init() {
	DefaultLogger, _ = New()
}

type (
	// Logger is the interface for all logging operations.
	Logger struct {
		level  Level
		logger log.Logger
	}

	// Level defines the logging levels available to the Logger, with a level of
	// debug logging all message and error logging only error message.
	Level uint

	// Format defines the output formats of the logger. Supported formats are
	// FormatLogfmt (the default), FormatJSON, and FormatNop (no logging).
	Format uint
)

const (
	// LevelDebug provides the most verbose logging, allowing logs for debug, info,
	// warn, and error.
	LevelDebug Level = iota

	// LevelInfo is the default log level. It allows logs for info, warn, and
	// error.
	LevelInfo

	// LevelWarn is reserved for logging warnings and errors.
	LevelWarn

	// LevelError provides the least verbose logging, allowing logs only for
	// errors.
	LevelError

	// LevelInvalid indicates and invalid log level.
	LevelInvalid

	// FormatJSON encodes each log entry as a single JSON object.
	FormatJSON Format = iota

	// FormatLogfmt encodes each log entry in logfmt format.
	FormatLogfmt

	// FormatNop will suppress logging output entirely.
	FormatNop

	// FormatInvalid indicates an invalid log format.
	FormatInvalid
)

// New returns a Logger.
func New(opts ...Option) (Logger, error) {
	o := options{
		level:       LevelInfo,
		format:      FormatLogfmt,
		writer:      os.Stdout,
		timestamped: true,
	}
	for _, opt := range opts {
		opt.apply(&o)
	}
	if o.level >= LevelInvalid {
		return Logger{}, fmt.Errorf("%w: %d", ErrInvalidLevel, o.level)
	}
	if o.format >= FormatInvalid {
		return Logger{}, fmt.Errorf("%w: %d", ErrInvalidFormat, o.format)
	}
	if o.writer == nil {
		return Logger{}, ErrInvalidOutput
	}

	l := Logger{level: o.level}

	switch o.format {
	case FormatJSON:
		l.logger = log.NewJSONLogger(log.NewSyncWriter(o.writer))
	case FormatLogfmt:
		l.logger = log.NewLogfmtLogger(log.NewSyncWriter(o.writer))
	case FormatNop:
		l.logger = log.NewNopLogger()
	}

	if o.timestamped {
		l.logger = log.With(l.logger, "ts", log.DefaultTimestampUTC)
	}
	if o.asDefault {
		DefaultLogger = l
	}

	return l, nil
}

// Debug logs a debug level message.
func (l Logger) Debug(kv ...interface{}) {
	if l.level > LevelDebug {
		return
	}
	level.Debug(l.logger).Log(kv...)
}

// Debug logs an debug level message to the default logger.
func Debug(kv ...interface{}) {
	DefaultLogger.Debug(kv...)
}

// Info logs an info level message.
func (l Logger) Info(kv ...interface{}) {
	if l.level > LevelInfo {
		return
	}
	level.Info(l.logger).Log(kv...)
}

// Info logs an info level message to the default logger.
func Info(kv ...interface{}) {
	DefaultLogger.Info(kv...)
}

// Warn logs a warn level message.
func (l Logger) Warn(kv ...interface{}) {
	if l.level > LevelWarn {
		return
	}
	level.Warn(l.logger).Log(kv...)
}

// Warn logs a warn level message to the default logger.
func Warn(kv ...interface{}) {
	DefaultLogger.Warn(kv...)
}

// Error logs an error level message.
func (l Logger) Error(kv ...interface{}) {
	level.Error(l.logger).Log(kv...)
}

// Error logs an error level message to the default logger.
func Error(kv ...interface{}) {
	DefaultLogger.Error(kv...)
}

// Level returns the log level.
func (l Logger) Level() Level {
	return l.level
}

// With returns a new contextual logger with keyvals prepended to those
// passed to calls to log.
func (l Logger) With(kv ...interface{}) Logger {
	return Logger{
		level:  l.level,
		logger: log.With(l.logger, kv...),
	}
}

// ToLevel translates the given level as a string to a Level.
func ToLevel(l string) Level {
	level := LevelInvalid
	switch strings.ToLower(l) {
	case "info", "": // An unset level string defaults to LevelInfo.
		level = LevelInfo
	case "debug":
		level = LevelDebug
	case "warn":
		level = LevelWarn
	case "error":
		level = LevelError
	default:
		level = LevelInvalid
	}
	return level
}

// ToFormat translates the given format as a string to a Format.
func ToFormat(f string) Format {
	format := FormatInvalid
	switch strings.ToLower(f) {
	case "logfmt", "": // An unset format string defaults to FormatLogfmt.
		format = FormatLogfmt
	case "json":
		format = FormatJSON
	case "nop":
		format = FormatNop
	default:
		format = FormatInvalid
	}
	return format
}

// NewContextWithLogger returns a new context with the given logger.
func NewContextWithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, logger)
}

// LoggerFromContext returns the logger for the current request.
func LoggerFromContext(ctx context.Context) Logger {
	logger, ok := ctx.Value(loggerContextKey).(Logger)
	if !ok {
		logger, _ = New(WithFormat(FormatNop))
	}
	return logger
}
