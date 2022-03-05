package log

import (
	"context"
	"errors"
	"os"
	"testing"

	"arcadium.dev/core/test"
)

func TestWithLevel(t *testing.T) {
	t.Parallel()

	levelCheck := func(t *testing.T, lvl Level) {
		t.Helper()
		var opts options
		o := WithLevel(lvl)
		if o == nil {
			t.Errorf("option expected")
		}
		o.apply(&opts)
		if opts.level != lvl {
			t.Errorf("Expected: %d actual: %d", lvl, opts.level)
		}
	}

	t.Run("LevelDebug", func(t *testing.T) {
		t.Parallel()
		levelCheck(t, LevelDebug)
	})

	t.Run("LevelInfo", func(t *testing.T) {
		t.Parallel()
		levelCheck(t, LevelDebug)
	})

	t.Run("LevelWarn", func(t *testing.T) {
		t.Parallel()
		levelCheck(t, LevelDebug)
	})

	t.Run("LevelError", func(t *testing.T) {
		t.Parallel()
		levelCheck(t, LevelDebug)
	})
}

func TestWithFormat(t *testing.T) {
	t.Parallel()

	formatCheck := func(t *testing.T, f Format) {
		t.Helper()
		var opts options
		o := WithFormat(f)
		if o == nil {
			t.Errorf("option expected")
		}
		o.apply(&opts)
		if opts.format != f {
			t.Errorf("Expected: %d actual: %d", f, opts.format)
		}
	}

	t.Run("FormatJSON", func(t *testing.T) {
		t.Parallel()
		formatCheck(t, FormatJSON)
	})

	t.Run("FormatLogfmt", func(t *testing.T) {
		t.Parallel()
		formatCheck(t, FormatJSON)
	})

	t.Run("FormatNop", func(t *testing.T) {
		t.Parallel()
		formatCheck(t, FormatJSON)
	})
}

func TestWithOutput(t *testing.T) {
	t.Parallel()
	var (
		opts options
		buf  = test.NewStringBuffer()
	)
	o := WithOutput(buf)
	if o == nil {
		t.Errorf("option expected")
	}
	o.apply(&opts)
	if opts.writer != buf {
		t.Errorf("Expected: %s actual: %s", buf, opts.writer)
	}
}

func TestWithoutTimestamp(t *testing.T) {
	t.Parallel()
	var opts options

	o := WithoutTimestamp()
	if o == nil {
		t.Errorf("option expected")
	}
	o.apply(&opts)
	if opts.timestamped != false {
		t.Error("Expected timestamped to be false")
	}
}

func TestNew(t *testing.T) {
	t.Parallel()

	t.Run("Test invalid level option", func(t *testing.T) {
		t.Parallel()
		l, err := New(WithLevel(Level(42)))
		if l != nil {
			t.Errorf("Unexpected logger")
		}
		if err == nil {
			t.Errorf("Expected an error")
		}
		if !errors.Is(err, ErrInvalidLevel) {
			t.Errorf("\nExpected: %s\nActual:   %s", ErrInvalidLevel, err)
		}
	})

	t.Run("Test invalid format option", func(t *testing.T) {
		t.Parallel()
		l, err := New(WithFormat(Format(42)))
		if l != nil {
			t.Errorf("Unexpected logger")
		}
		if err == nil {
			t.Errorf("Expected an error")
		}
		if !errors.Is(err, ErrInvalidFormat) {
			t.Errorf("\nExpected: %s\nActual:   %s", ErrInvalidFormat, err)
		}
	})
	t.Run("Test invalid output option", func(t *testing.T) {
		t.Parallel()
		l, err := New(WithOutput(nil))
		if l != nil {
			t.Errorf("Unexpected logger")
		}
		if err == nil {
			t.Errorf("Expected an error")
		}
		if !errors.Is(err, ErrInvalidOutput) {
			t.Errorf("\nExpected: %s\nActual:   %s", ErrInvalidOutput, err)
		}
	})

	t.Run("Test valid level option", func(t *testing.T) {
		t.Parallel()
		l, err := New(WithLevel(LevelError))
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil || l.(*logger).logger == nil {
			t.Error("Expected a logger")
		}
	})

	t.Run("Test valid format option", func(t *testing.T) {
		t.Parallel()
		l, err := New(WithFormat(FormatNop))
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil || l.(*logger).logger == nil {
			t.Error("Expected a logger")
		}
	})

	t.Run("Test valid output option", func(t *testing.T) {
		t.Parallel()
		l, err := New(WithOutput(test.NewStringBuffer()))
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil || l.(*logger).logger == nil {
			t.Error("Expected a logger")
		}
	})

	t.Run("Test all options", func(t *testing.T) {
		t.Parallel()
		l, err := New(
			WithLevel(LevelWarn),
			WithFormat(FormatLogfmt),
			WithOutput(os.Stdout),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil || l.(*logger).logger == nil {
			t.Error("Expected a logger")
		}
	})

	t.Run("Test defaults ", func(t *testing.T) {
		t.Parallel()
		lg, err := New()
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if lg == nil {
			t.Errorf("Expected a logger")
		}
		l, ok := lg.(*logger)
		if !ok {
			t.Errorf("Failed logger type assertion")
		}

		if l.opts.level != LevelInfo {
			t.Errorf("Expect %d: Actual %d", LevelInfo, l.opts.level)
		}
		if l.opts.format != FormatJSON {
			t.Errorf("Expect %d: Actual %d", FormatJSON, l.opts.format)
		}
		if l.opts.writer != os.Stderr {
			t.Errorf("Unexpected writer: %+v", l.opts.writer)
		}
	})
}

func TestDebug(t *testing.T) {
	t.Parallel()

	t.Run("Level set to debug", func(t *testing.T) {
		t.Parallel()
		b := test.NewStringBuffer()
		l, err := New(
			WithLevel(LevelDebug),
			WithFormat(FormatLogfmt),
			WithOutput(b),
			WithoutTimestamp(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil {
			t.Errorf("Expected a logger")
		}

		l.Debug("a", "b")
		l.Debug("c", "d")

		if len(b.Buffer) != 2 {
			t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
		}
		expected := "level=debug a=b\n"
		if b.Buffer[0] != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Buffer[0])
		}
		expected = "level=debug c=d\n"
		if b.Buffer[1] != expected {
			t.Errorf("Expected '%s', Actual: '%s'", expected, b.Buffer[1])
		}
	})

	t.Run("Level set above debug", func(t *testing.T) {
		t.Parallel()
		b := test.NewStringBuffer()
		l, err := New(
			WithLevel(LevelInfo),
			WithFormat(FormatLogfmt),
			WithOutput(b),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil {
			t.Errorf("Expected a logger")
		}

		l.Debug("a", "b")

		if len(b.Buffer) != 0 {
			t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
		}
	})
}

func TestInfo(t *testing.T) {
	t.Parallel()

	t.Run("Level set to info", func(t *testing.T) {
		t.Parallel()
		b := test.NewStringBuffer()
		l, err := New(
			WithLevel(LevelInfo),
			WithFormat(FormatJSON),
			WithOutput(b),
			WithoutTimestamp(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil {
			t.Errorf("Expected a logger")
		}

		l.Info("a", "b")
		l.Info("c", "d")

		if len(b.Buffer) != 2 {
			t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
		}
		expected := "{\"a\":\"b\",\"level\":\"info\"}\n"
		if b.Buffer[0] != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Buffer[0])
		}
		expected = "{\"c\":\"d\",\"level\":\"info\"}\n"
		if b.Buffer[1] != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Buffer[1])
		}
	})

	t.Run("Level set above info", func(t *testing.T) {
		t.Parallel()
		b := test.NewStringBuffer()
		l, err := New(
			WithLevel(LevelWarn),
			WithFormat(FormatJSON),
			WithOutput(b),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil {
			t.Errorf("Expected a logger")
		}

		l.Info("a", "b")

		if len(b.Buffer) != 0 {
			t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
		}
	})
}

func TestWarn(t *testing.T) {
	t.Parallel()

	t.Run("Level set to warn", func(t *testing.T) {
		t.Parallel()
		b := test.NewStringBuffer()
		l, err := New(
			WithLevel(LevelWarn),
			WithFormat(FormatLogfmt),
			WithOutput(b),
			WithoutTimestamp(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil {
			t.Errorf("Expected a logger")
		}

		l.Warn("a", "b")
		l.Warn("c", "d")

		if len(b.Buffer) != 2 {
			t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
		}
		expected := "level=warn a=b\n"
		if b.Buffer[0] != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Buffer[0])
		}
		expected = "level=warn c=d\n"
		if b.Buffer[1] != expected {
			t.Errorf("Expected '%s', Actual: '%s'", expected, b.Buffer[1])
		}
	})

	t.Run("Level set above warn", func(t *testing.T) {
		t.Parallel()
		b := test.NewStringBuffer()
		l, err := New(
			WithLevel(LevelError),
			WithFormat(FormatLogfmt),
			WithOutput(b),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil {
			t.Errorf("Expected a logger")
		}

		l.Warn("a", "b")

		if len(b.Buffer) != 0 {
			t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
		}
	})
}

func TestError(t *testing.T) {
	t.Parallel()

	t.Run("Level set to error", func(t *testing.T) {
		t.Parallel()
		b := test.NewStringBuffer()
		l, err := New(
			WithLevel(LevelError),
			WithFormat(FormatJSON),
			WithOutput(b),
			WithoutTimestamp(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil {
			t.Errorf("Expected a logger")
		}

		l.Error("a", "b")
		l.Error("c", "d")

		if len(b.Buffer) != 2 {
			t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
		}
		expected := "{\"a\":\"b\",\"level\":\"error\"}\n"
		if b.Buffer[0] != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Buffer[0])
		}
		expected = "{\"c\":\"d\",\"level\":\"error\"}\n"
		if b.Buffer[1] != expected {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Buffer[1])
		}
	})
}

func TestLogging(t *testing.T) {
	t.Parallel()
	b := test.NewStringBuffer()
	l, err := New(
		WithLevel(LevelDebug),
		WithOutput(b),
		WithoutTimestamp(),
	)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if l == nil {
		t.Errorf("Expected a logger")
	}

	l.Debug("a", "b", "c", "d", "e", "f")
	l.Info("foo", "bar", "beh", "baz")
	l.Warn("msg", "Hello World", "addr", "127.0.0.1", "port", 8443)
	l.Error("msg", "this is an error")

	if len(b.Buffer) != 4 {
		t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
	}
	expected := "{\"a\":\"b\",\"c\":\"d\",\"e\":\"f\",\"level\":\"debug\"}\n"
	if b.Buffer[0] != expected {
		t.Errorf("\nExpected %sActual:  %s", expected, b.Buffer[0])
	}
	expected = "{\"beh\":\"baz\",\"foo\":\"bar\",\"level\":\"info\"}\n"
	if b.Buffer[1] != expected {
		t.Errorf("\nExpected %sActual:  %s", expected, b.Buffer[1])
	}
	expected = "{\"addr\":\"127.0.0.1\",\"level\":\"warn\",\"msg\":\"Hello World\",\"port\":8443}\n"
	if b.Buffer[2] != expected {
		t.Errorf("\nExpected %sActual:  %s", expected, b.Buffer[2])
	}
	expected = "{\"level\":\"error\",\"msg\":\"this is an error\"}\n"
	if b.Buffer[3] != expected {
		t.Errorf("\nExpected %sActual:  %s", expected, b.Buffer[3])
	}
}

func TestWith(t *testing.T) {
	t.Parallel()
	b := test.NewStringBuffer()
	l, err := New(
		WithFormat(FormatLogfmt),
		WithOutput(b),
		WithoutTimestamp(),
	)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if l == nil {
		t.Errorf("Expected a logger")
	}
	l = l.With("id", "0000-111")

	l.Info("bob", "smith", "alice", "jones")

	if len(b.Buffer) != 1 {
		t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
	}

	expected := "level=info id=0000-111 bob=smith alice=jones\n"
	if b.Buffer[0] != expected {
		t.Errorf("\nExpected %sActual:  %s", expected, b.Buffer[0])
	}
}

func TestToLevel(t *testing.T) {
	levels := []struct {
		s string
		l Level
	}{
		{s: "", l: LevelInfo},
		{s: "info", l: LevelInfo},
		{s: "debug", l: LevelDebug},
		{s: "DEBUG", l: LevelDebug},
		{s: "INFO", l: LevelInfo},
		{s: "warn", l: LevelWarn},
		{s: "WARN", l: LevelWarn},
		{s: "error", l: LevelError},
		{s: "ERROR", l: LevelError},
		{s: "invalid", l: LevelInvalid},
	}
	for _, l := range levels {
		if ToLevel(l.s) != l.l {
			t.Errorf("Unexpected level: %s, for %d", l.s, l.l)
		}
	}
}

func TestToFormat(t *testing.T) {
	formats := []struct {
		s string
		f Format
	}{
		{s: "", f: FormatJSON},
		{s: "json", f: FormatJSON},
		{s: "JSON", f: FormatJSON},
		{s: "logfmt", f: FormatLogfmt},
		{s: "LOGFMT", f: FormatLogfmt},
		{s: "nop", f: FormatNop},
		{s: "NOP", f: FormatNop},
		{s: "invalid", f: FormatInvalid},
	}
	for _, f := range formats {
		if ToFormat(f.s) != f.f {
			t.Errorf("Unexpected level: %s, for %d", f.s, f.f)
		}
	}
}

func TestLoggerContext(t *testing.T) {
	t.Parallel()
	b := test.NewStringBuffer()
	l, err := New(
		WithFormat(FormatLogfmt),
		WithOutput(b),
		WithoutTimestamp(),
	)
	if err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
	if l == nil {
		t.Errorf("Expected a logger")
	}
	l = l.With("id", "0000-111")

	ctx := context.Background()

	ctx = NewContextWithLogger(ctx, l)
	nl := LoggerFromContext(ctx)

	if l != nl {
		t.Errorf("Unexpected logger from context")
	}

	nl.Info("bob", "smith", "alice", "jones")

	if len(b.Buffer) != 1 {
		t.Errorf("Unexpected buffer length: %d", len(b.Buffer))
	}

	expected := "level=info id=0000-111 bob=smith alice=jones\n"
	if b.Buffer[0] != expected {
		t.Errorf("\nExpected %sActual:  %s", expected, b.Buffer[0])
	}
}
