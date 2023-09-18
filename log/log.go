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

// Package log provides tools to setup the zerolog logger.
package log // import "arcadium.dev/core/log"

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var (
	// ErrInvalidLevel will be returned when the level given to the WithLevel
	// option is invalid.
	ErrInvalidLevel = errors.New("invalid level")

	// ErrInvalidOutput will be returned when the output writer given to WithOutput
	// is nil.
	ErrInvalidOutput = errors.New("invalid output")
)

// New returns a zerolog.Logger.
func New(opts ...Option) (*zerolog.Logger, error) {
	o := options{
		level:              zerolog.InfoLevel,
		levelFieldName:     zerolog.LevelFieldName,
		timestamped:        true,
		timestampFieldName: zerolog.TimestampFieldName,
		messageFieldName:   zerolog.MessageFieldName,
		writer:             os.Stderr,
	}
	for _, opt := range opts {
		opt.apply(&o)
	}
	if o.level >= zerolog.NoLevel {
		return nil, fmt.Errorf("%w: %d", ErrInvalidLevel, o.level)
	}
	if o.writer == nil {
		return nil, ErrInvalidOutput
	}

	nop := zerolog.Nop()
	zerolog.DefaultContextLogger = &nop

	zerolog.LevelFieldName = o.levelFieldName
	zerolog.SetGlobalLevel(o.level)

	zerolog.MessageFieldName = o.messageFieldName

	zerolog.TimestampFieldName = o.timestampFieldName
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.TimestampFunc = func() time.Time {
		return time.Now().UTC()
	}

	l := zerolog.New(o.writer)

	if o.timestamped {
		l = l.With().Timestamp().Logger()
	}
	if o.asDefault {
		log.Logger = l
	}

	return &l, nil
}

// ToLevel translates the given level as a string to a Level.
func ToLevel(l string) zerolog.Level {
	level := zerolog.NoLevel
	switch strings.ToLower(l) {
	case "info", "": // An unset level string defaults to LevelInfo.
		level = zerolog.InfoLevel
	case "debug":
		level = zerolog.DebugLevel
	case "warn":
		level = zerolog.WarnLevel
	case "error":
		level = zerolog.ErrorLevel
	}
	return level
}
