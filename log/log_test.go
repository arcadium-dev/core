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

package log_test

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"arcadium.dev/core/log"
)

func TestNew(t *testing.T) {
	t.Run("Test invalid level option", func(t *testing.T) {
		_, err := log.New(log.WithLevel(log.Level(42)))
		if err == nil {
			t.Errorf("Expected an error")
		}
		if !errors.Is(err, log.ErrInvalidLevel) {
			t.Errorf("\nExpected: %s\nActual:   %s", log.ErrInvalidLevel, err)
		}
	})

	t.Run("Test invalid format option", func(t *testing.T) {
		_, err := log.New(log.WithFormat(log.Format(42)))
		if err == nil {
			t.Errorf("Expected an error")
		}
		if !errors.Is(err, log.ErrInvalidFormat) {
			t.Errorf("\nExpected: %s\nActual:   %s", log.ErrInvalidFormat, err)
		}
	})
	t.Run("Test invalid output option", func(t *testing.T) {
		_, err := log.New(log.WithOutput(nil))
		if err == nil {
			t.Errorf("Expected an error")
		}
		if !errors.Is(err, log.ErrInvalidOutput) {
			t.Errorf("\nExpected: %s\nActual:   %s", log.ErrInvalidOutput, err)
		}
	})

	t.Run("Test valid level option", func(t *testing.T) {
		l, err := log.New(log.WithLevel(log.LevelError))
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.LevelError != l.Level() {
			t.Errorf("Expected level: %d\nActual level:   %d", log.LevelError, l.Level())
		}
	})

	t.Run("Test valid format option", func(t *testing.T) {
		_, err := log.New(log.WithFormat(log.FormatNop))
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
	})

	t.Run("Test valid output option", func(t *testing.T) {
		_, err := log.New(log.WithOutput(log.NewStringBuffer()))
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
	})

	t.Run("Test as default option", func(t *testing.T) {
		l, err := log.New(log.AsDefault())
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.DefaultLogger != l {
			t.Errorf("Expected DefaultLogger to be the same as the new logger")
		}
	})

	t.Run("Test all options", func(t *testing.T) {
		l, err := log.New(
			log.WithLevel(log.LevelWarn),
			log.WithFormat(log.FormatLogfmt),
			log.WithOutput(os.Stdout),
			log.AsDefault(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.DefaultLogger != l {
			t.Errorf("Expected DefaultLogger to be the same as the new logger")
		}
	})

	t.Run("Test defaults ", func(t *testing.T) {
		l, err := log.New()
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}

		if log.LevelInfo != l.Level() {
			t.Errorf("Expect level: %d\nActual level:   %d", log.LevelInfo, l.Level())
		}
	})
}

func TestDebug(t *testing.T) {
	t.Run("Level set to debug", func(t *testing.T) {
		b := log.NewStringBuffer()
		l, err := log.New(
			log.WithLevel(log.LevelDebug),
			log.WithFormat(log.FormatLogfmt),
			log.WithOutput(b),
			log.WithoutTimestamp(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.LevelDebug != l.Level() {
			t.Errorf("Expect level: %d\nActual level:   %d", log.LevelDebug, l.Level())
		}

		l.Debug("a", "b")
		l.Debug("c", "d")

		if b.Len() != 2 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
		expected := "level=debug a=b\n"
		if b.Index(0) != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Index(0))
		}
		expected = "level=debug c=d\n"
		if b.Index(1) != expected {
			t.Errorf("Expected '%s', Actual: '%s'", expected, b.Index(1))
		}
	})

	t.Run("Level set above debug", func(t *testing.T) {
		b := log.NewStringBuffer()
		l, err := log.New(
			log.WithLevel(log.LevelInfo),
			log.WithFormat(log.FormatLogfmt),
			log.WithOutput(b),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.LevelInfo != l.Level() {
			t.Errorf("Expect level: %d\nActual level:   %d", log.LevelInfo, l.Level())
		}

		l.Debug("a", "b")

		if b.Len() != 0 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
	})

	t.Run("Test global debug", func(t *testing.T) {
		b := log.NewStringBuffer()
		l, err := log.New(
			log.WithLevel(log.LevelDebug),
			log.WithOutput(b),
			log.WithoutTimestamp(),
			log.AsDefault(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.LevelDebug != l.Level() {
			t.Errorf("Expect level: %d\nActual level:   %d", log.LevelDebug, l.Level())
		}
		if log.DefaultLogger != l {
			t.Errorf("Expected default logger to be set")
		}

		log.Debug("a", "b")
		log.Debug("c", "d")

		if b.Len() != 2 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
		expected := "level=debug a=b\n"
		if b.Index(0) != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Index(0))
		}
		expected = "level=debug c=d\n"
		if b.Index(1) != expected {
			t.Errorf("Expected '%s', Actual: '%s'", expected, b.Index(1))
		}
	})
}

func TestInfo(t *testing.T) {
	t.Run("Level set to info", func(t *testing.T) {
		b := log.NewStringBuffer()
		l, err := log.New(
			log.WithLevel(log.LevelInfo),
			log.WithFormat(log.FormatJSON),
			log.WithOutput(b),
			log.WithoutTimestamp(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.LevelInfo != l.Level() {
			t.Errorf("Expected level: %d\nActual level:   %d", log.LevelInfo, l.Level())
		}

		l.Info("a", "b")
		l.Info("c", "d")

		if b.Len() != 2 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
		expected := "{\"a\":\"b\",\"level\":\"info\"}\n"
		if b.Index(0) != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Index(0))
		}
		expected = "{\"c\":\"d\",\"level\":\"info\"}\n"
		if b.Index(1) != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Index(1))
		}
	})

	t.Run("Level set above info", func(t *testing.T) {
		b := log.NewStringBuffer()
		l, err := log.New(
			log.WithLevel(log.LevelWarn),
			log.WithFormat(log.FormatJSON),
			log.WithOutput(b),
			log.WithoutTimestamp(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.LevelWarn != l.Level() {
			t.Errorf("Expected level: %d\nActual level:   %d", log.LevelWarn, l.Level())
		}

		l.Info("a", "b")

		if b.Len() != 0 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
	})

	t.Run("Test global info", func(t *testing.T) {
		b := log.NewStringBuffer()
		l, err := log.New(
			log.WithLevel(log.LevelInfo),
			log.WithOutput(b),
			log.WithoutTimestamp(),
			log.AsDefault(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.LevelInfo != l.Level() {
			t.Errorf("Expected level: %d\nActual level:   %d", log.LevelInfo, l.Level())
		}
		if log.DefaultLogger != l {
			t.Errorf("Expected default logger to be set")
		}

		log.Info("a", "b")
		log.Info("c", "d")

		if b.Len() != 2 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
		expected := "level=info a=b\n"
		if b.Index(0) != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Index(0))
		}
		expected = "level=info c=d\n"
		if b.Index(1) != expected {
			t.Errorf("Expected '%s', Actual: '%s'", expected, b.Index(1))
		}
	})
}

func TestWarn(t *testing.T) {
	t.Run("Level set to warn", func(t *testing.T) {
		b := log.NewStringBuffer()
		l, err := log.New(
			log.WithLevel(log.LevelWarn),
			log.WithFormat(log.FormatLogfmt),
			log.WithOutput(b),
			log.WithoutTimestamp(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.LevelWarn != l.Level() {
			t.Errorf("Expected level: %d\nActual level:   %d", log.LevelWarn, l.Level())
		}

		l.Warn("a", "b")
		l.Warn("c", "d")

		if b.Len() != 2 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
		expected := "level=warn a=b\n"
		if b.Index(0) != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Index(0))
		}
		expected = "level=warn c=d\n"
		if b.Index(1) != expected {
			t.Errorf("Expected '%s', Actual: '%s'", expected, b.Index(1))
		}
	})

	t.Run("Level set above warn", func(t *testing.T) {
		b := log.NewStringBuffer()
		l, err := log.New(
			log.WithLevel(log.LevelError),
			log.WithFormat(log.FormatLogfmt),
			log.WithOutput(b),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.LevelError != l.Level() {
			t.Errorf("Expected level: %d\nActual level:   %d", log.LevelError, l.Level())
		}

		l.Warn("a", "b")

		if b.Len() != 0 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
	})

	t.Run("Test global warn", func(t *testing.T) {
		b := log.NewStringBuffer()
		l, err := log.New(
			log.WithLevel(log.LevelWarn),
			log.WithOutput(b),
			log.WithoutTimestamp(),
			log.AsDefault(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.LevelWarn != l.Level() {
			t.Errorf("Expected level: %d\nActual level:   %d", log.LevelWarn, l.Level())
		}
		if log.DefaultLogger != l {
			t.Errorf("Expected default logger to be set")
		}

		log.Warn("a", "b")
		log.Warn("c", "d")

		if b.Len() != 2 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
		expected := "level=warn a=b\n"
		if b.Index(0) != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Index(0))
		}
		expected = "level=warn c=d\n"
		if b.Index(1) != expected {
			t.Errorf("Expected '%s', Actual: '%s'", expected, b.Index(1))
		}
	})
}

func TestError(t *testing.T) {
	t.Run("Level set to error", func(t *testing.T) {
		b := log.NewStringBuffer()
		l, err := log.New(
			log.WithLevel(log.LevelError),
			log.WithFormat(log.FormatJSON),
			log.WithOutput(b),
			log.WithoutTimestamp(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.LevelError != l.Level() {
			t.Errorf("Expected level: %d\nActual level:   %d", log.LevelError, l.Level())
		}

		l.Error("a", "b")
		l.Error("c", "d")

		if b.Len() != 2 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
		expected := "{\"a\":\"b\",\"level\":\"error\"}\n"
		if b.Index(0) != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Index(0))
		}
		expected = "{\"c\":\"d\",\"level\":\"error\"}\n"
		if b.Index(1) != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Index(1))
		}
	})

	t.Run("Test global error", func(t *testing.T) {
		b := log.NewStringBuffer()
		l, err := log.New(
			log.WithLevel(log.LevelError),
			log.WithOutput(b),
			log.WithoutTimestamp(),
			log.AsDefault(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if log.LevelError != l.Level() {
			t.Errorf("Expected level: %d\nActual level:   %d", log.LevelError, l.Level())
		}
		if log.DefaultLogger != l {
			t.Errorf("Expected default logger to be set")
		}

		log.Error("a", "b")
		log.Error("c", "d")

		if b.Len() != 2 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
		expected := "level=error a=b\n"
		if b.Index(0) != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Index(0))
		}
		expected = "level=error c=d\n"
		if b.Index(1) != expected {
			t.Errorf("Expected '%s', Actual: '%s'", expected, b.Index(1))
		}
	})
}

func TestLogging(t *testing.T) {
	b := log.NewStringBuffer()
	l, err := log.New(
		log.WithLevel(log.LevelDebug),
		log.WithOutput(b),
		log.WithoutTimestamp(),
	)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}

	l.Debug("a", "b", "c", "d", "e", "f")
	l.Info("foo", "bar", "beh", "baz")
	l.Warn("msg", "Hello World", "addr", "127.0.0.1", "port", 8443)
	l.Error("msg", "this is an error")

	if b.Len() != 4 {
		t.Errorf("Unexpected buffer length: %d", b.Len())
	}
	expected := "level=debug a=b c=d e=f\n"
	if b.Index(0) != expected {
		t.Errorf("\nExpected %s\nActual:  %s", expected, b.Index(0))
	}
	expected = "level=info foo=bar beh=baz\n"
	if b.Index(1) != expected {
		t.Errorf("\nExpected %s\nActual:  %s", expected, b.Index(1))
	}
	expected = "level=warn msg=\"Hello World\" addr=127.0.0.1 port=8443\n"
	if b.Index(2) != expected {
		t.Errorf("\nExpected %s\nActual:  %s", expected, b.Index(2))
	}
	expected = "level=error msg=\"this is an error\"\n"
	if b.Index(3) != expected {
		t.Errorf("\nExpected %s\nActual:  %s", expected, b.Index(3))
	}
}

func TestWith(t *testing.T) {
	b := log.NewStringBuffer()
	l, err := log.New(
		log.WithFormat(log.FormatLogfmt),
		log.WithOutput(b),
		log.WithoutTimestamp(),
	)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	l = l.With("id", "0000-111")

	l.Info("bob", "smith", "alice", "jones")

	if b.Len() != 1 {
		t.Errorf("Unexpected buffer length: %d", b.Len())
	}

	expected := "level=info id=0000-111 bob=smith alice=jones\n"
	if b.Index(0) != expected {
		t.Errorf("\nExpected %sActual:  %s", expected, b.Index(0))
	}
}

func TestToLevel(t *testing.T) {
	levels := []struct {
		s string
		l log.Level
	}{
		{s: "", l: log.LevelInfo},
		{s: "info", l: log.LevelInfo},
		{s: "debug", l: log.LevelDebug},
		{s: "DEBUG", l: log.LevelDebug},
		{s: "INFO", l: log.LevelInfo},
		{s: "warn", l: log.LevelWarn},
		{s: "WARN", l: log.LevelWarn},
		{s: "error", l: log.LevelError},
		{s: "ERROR", l: log.LevelError},
		{s: "invalid", l: log.LevelInvalid},
	}
	for _, l := range levels {
		if log.ToLevel(l.s) != l.l {
			t.Errorf("Unexpected level: %s, for %d", l.s, l.l)
		}
	}
}

func TestToFormat(t *testing.T) {
	formats := []struct {
		s string
		f log.Format
	}{
		{s: "", f: log.FormatLogfmt},
		{s: "json", f: log.FormatJSON},
		{s: "JSON", f: log.FormatJSON},
		{s: "logfmt", f: log.FormatLogfmt},
		{s: "LOGFMT", f: log.FormatLogfmt},
		{s: "nop", f: log.FormatNop},
		{s: "NOP", f: log.FormatNop},
		{s: "invalid", f: log.FormatInvalid},
	}
	for _, f := range formats {
		if log.ToFormat(f.s) != f.f {
			t.Errorf("Unexpected level: %s, for %d", f.s, f.f)
		}
	}
}

func TestLoggerContext(t *testing.T) {
	t.Run("insert, extract success", func(t *testing.T) {
		b := log.NewStringBuffer()
		l, err := log.New(
			log.WithFormat(log.FormatLogfmt),
			log.WithOutput(b),
			log.WithoutTimestamp(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		l = l.With("id", "0000-111")

		ctx := context.Background()

		ctx = log.NewContextWithLogger(ctx, l)
		nl := log.LoggerFromContext(ctx)

		if l != nl {
			t.Errorf("Unexpected logger from context")
		}

		nl.Info("bob", "smith", "alice", "jones")

		if b.Len() != 1 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}

		expected := "level=info id=0000-111 bob=smith alice=jones\n"
		if !strings.Contains(b.Index(0), expected) {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Index(0))
		}
	})
}
