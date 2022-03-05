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

package log // import "arcadium.dev/core/log"

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

type (
	// Logger is the interface for all logging operations.
	Logger interface {
		// Debug logs a debug level message.
		Debug(...interface{})

		// Info logs an info level message.
		Info(...interface{})

		// Warn logs a warn level message.
		Warn(...interface{})

		// Error logs an error level message.
		Error(...interface{})

		// With returns a new contextual logger with keyvals prepended to those
		// passed to calls to log.
		With(...interface{}) Logger
	}

	// Option provides for Logger configuration.
	Option interface {
		apply(*options)
	}

	// Level defines the logging levels available to the Logger, with a level of
	// debug logging all message and error logging only error message.
	Level uint

	// Format defines the output formats of the logger. Supported formats are
	// FormatJSON (the default), FormatLogfmt, and FormatNop (no loggin).
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

var (
	// ErrInvalidLevel will be returned when the level given to the WirhLevel
	// option is invalid.
	ErrInvalidLevel = errors.New("Invalid Level")

	// ErrInvalidFormat will be returned when the format given to the WithFormat
	// options is invalid.
	ErrInvalidFormat = errors.New("Invvalid Format")

	// ErrInvalidOutput will be returned when the output writer given to WithOuput
	// is nil.
	ErrInvalidOutput = errors.New("Invalid Format")
)

// New returns a Logger.
func New(opts ...Option) (Logger, error) {
	l := &logger{
		opts: options{
			level:       LevelInfo,
			format:      FormatJSON,
			writer:      os.Stderr,
			timestamped: true,
		},
	}

	for _, opt := range opts {
		opt.apply(&l.opts)
	}
	if l.opts.level >= LevelInvalid {
		return nil, fmt.Errorf("%w: %d", ErrInvalidLevel, l.opts.level)
	}
	if l.opts.format >= FormatInvalid {
		return nil, fmt.Errorf("%w: %d", ErrInvalidFormat, l.opts.format)
	}
	if l.opts.writer == nil {
		return nil, ErrInvalidOutput
	}

	switch l.opts.format {
	case FormatJSON:
		l.logger = log.NewJSONLogger(log.NewSyncWriter(l.opts.writer))
	case FormatLogfmt:
		l.logger = log.NewLogfmtLogger(log.NewSyncWriter(l.opts.writer))
	case FormatNop:
		l.logger = log.NewNopLogger()
	}
	if l.opts.timestamped {
		l.logger = log.With(l.logger, "ts", log.DefaultTimestampUTC)
	}

	return l, nil
}

// WithLevel allows the level to be configured. The default level is LevelInfo.
func WithLevel(level Level) Option {
	return newOption(func(opts *options) {
		opts.level = level
	})
}

// WithFormat allows the format to be configured. The default format is
// FormatJSON.
func WithFormat(format Format) Option {
	return newOption(func(opts *options) {
		opts.format = format
	})
}

// WithOutput allows the format to be configured. The default writer is
// os.Stderr.
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

// Debug logs a debug level message.
func (l *logger) Debug(kv ...interface{}) {
	if l.opts.level > LevelDebug {
		return
	}
	level.Debug(l.logger).Log(kv...)
}

// Info logs an info level message.
func (l *logger) Info(kv ...interface{}) {
	if l.opts.level > LevelInfo {
		return
	}
	level.Info(l.logger).Log(kv...)
}

// Warn logs a warn level message.
func (l *logger) Warn(kv ...interface{}) {
	if l.opts.level > LevelWarn {
		return
	}
	level.Warn(l.logger).Log(kv...)
}

// Error logs an error level message.
func (l *logger) Error(kv ...interface{}) {
	level.Error(l.logger).Log(kv...)
}

// With returns a new contextual logger with keyvals prepended to those
// passed to calls to log.
func (l *logger) With(kv ...interface{}) Logger {
	return &logger{
		opts:   l.opts,
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
	case "json", "": // An unset format string defaults to FormatJSON.
		format = FormatJSON
	case "logfmt":
		format = FormatLogfmt
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
	logger, _ := ctx.Value(loggerContextKey).(Logger)
	return logger
}

type (
	logger struct {
		opts   options
		logger log.Logger
	}

	options struct {
		level       Level
		format      Format
		writer      io.Writer
		timestamped bool
	}

	option struct {
		f func(*options)
	}

	contextKey int
)

const (
	loggerContextKey = contextKey(iota + 1)
)

func newOption(f func(*options)) *option {
	return &option{f: f}
}

func (o *option) apply(opts *options) {
	o.f(opts)
}
