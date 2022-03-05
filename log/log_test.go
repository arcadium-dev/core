package log

import (
	"context"
	"errors"
	"os"
	"strings"
	"testing"

	"arcadium.dev/core/test"
)

func TestWithLevel(t *testing.T) {
	levelCheck := func(t *testing.T, lvl Level) {
		t.Helper()
		var opts options
		o := WithLevel(lvl)
		if o == nil {
			t.Fatal("option expected")
		}
		o.apply(&opts)
		if opts.level != lvl {
			t.Errorf("Expected: %d actual: %d", lvl, opts.level)
		}
	}

	t.Run("LevelDebug", func(t *testing.T) {
		levelCheck(t, LevelDebug)
	})

	t.Run("LevelInfo", func(t *testing.T) {
		levelCheck(t, LevelDebug)
	})

	t.Run("LevelWarn", func(t *testing.T) {
		levelCheck(t, LevelDebug)
	})

	t.Run("LevelError", func(t *testing.T) {
		levelCheck(t, LevelDebug)
	})
}

func TestWithFormat(t *testing.T) {
	formatCheck := func(t *testing.T, f Format) {
		t.Helper()
		var opts options
		o := WithFormat(f)
		if o == nil {
			t.Fatal("option expected")
		}
		o.apply(&opts)
		if opts.format != f {
			t.Errorf("Expected: %d actual: %d", f, opts.format)
		}
	}

	t.Run("FormatJSON", func(t *testing.T) {
		formatCheck(t, FormatJSON)
	})

	t.Run("FormatLogfmt", func(t *testing.T) {
		formatCheck(t, FormatJSON)
	})

	t.Run("FormatNop", func(t *testing.T) {
		formatCheck(t, FormatJSON)
	})
}

func TestWithOutput(t *testing.T) {
	var (
		opts options
		b    = test.NewStringBuffer()
	)
	o := WithOutput(b)
	if o == nil {
		t.Fatal("option expected")
	}
	o.apply(&opts)
	if opts.writer != b {
		t.Errorf("Expected: %+v actual: %+v", b, opts.writer)
	}
}

func TestWithoutTimestamp(t *testing.T) {
	var opts options

	o := WithoutTimestamp()
	if o == nil {
		t.Fatal("option expected")
	}
	o.apply(&opts)
	if opts.timestamped != false {
		t.Error("Expected timestamped to be false")
	}
}

func TestAsDefault(t *testing.T) {
	var opts options

	o := AsDefault()
	if o == nil {
		t.Fatal("option expected")
	}
	o.apply(&opts)
	if !opts.asDefault {
		t.Error("Expected asDefault to be true")
	}
}

func TestNew(t *testing.T) {
	t.Run("Test invalid level option", func(t *testing.T) {
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
		l, err := New(WithLevel(LevelError))
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil || l.(*logger).logger == nil {
			t.Error("Expected a logger")
		}
	})

	t.Run("Test valid format option", func(t *testing.T) {
		l, err := New(WithFormat(FormatNop))
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil || l.(*logger).logger == nil {
			t.Error("Expected a logger")
		}
	})

	t.Run("Test valid output option", func(t *testing.T) {
		l, err := New(WithOutput(test.NewStringBuffer()))
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil || l.(*logger).logger == nil {
			t.Error("Expected a logger")
		}
	})

	t.Run("Test as default option", func(t *testing.T) {
		l, err := New(AsDefault())
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil || DefaultLogger != l {
			t.Errorf("Expected DefaultLogger to be the same as the new logger")
		}
	})

	t.Run("Test all options", func(t *testing.T) {
		l, err := New(
			WithLevel(LevelWarn),
			WithFormat(FormatLogfmt),
			WithOutput(os.Stdout),
			AsDefault(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if l == nil || l.(*logger).logger == nil {
			t.Error("Expected a logger")
		}
		if DefaultLogger != l {
			t.Errorf("Expected DefaultLogger to be the same as the new logger")
		}
	})

	t.Run("Test defaults ", func(t *testing.T) {
		DefaultLogger, err := New()
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		l, ok := DefaultLogger.(*logger)
		if !ok || l == nil {
			t.Fatalf("Failed to initialize default logger")
		}

		if l.opts.level != LevelInfo {
			t.Errorf("Expect %d: Actual %d", LevelInfo, l.opts.level)
		}
		if l.opts.format != FormatLogfmt {
			t.Errorf("Expect %d: Actual %d", FormatLogfmt, l.opts.format)
		}
		if l.opts.writer != os.Stdout {
			t.Errorf("Unexpected writer: %+v", l.opts.writer)
		}
	})
}

func TestDebug(t *testing.T) {
	t.Run("Level set to debug", func(t *testing.T) {
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

		if b.Len() != 0 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
	})

	t.Run("Test global debug", func(t *testing.T) {
		b := test.NewStringBuffer()
		l, err := New(
			WithLevel(LevelDebug),
			WithOutput(b),
			WithoutTimestamp(),
			AsDefault(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if DefaultLogger != l {
			t.Errorf("Expected default logger to be set")
		}

		Debug("a", "b")
		Debug("c", "d")

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
		b := test.NewStringBuffer()
		l, err := New(
			WithLevel(LevelWarn),
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

		if b.Len() != 0 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
	})

	t.Run("Test global info", func(t *testing.T) {
		b := test.NewStringBuffer()
		l, err := New(
			WithLevel(LevelInfo),
			WithOutput(b),
			WithoutTimestamp(),
			AsDefault(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if DefaultLogger != l {
			t.Errorf("Expected default logger to be set")
		}

		Info("a", "b")
		Info("c", "d")

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

		if b.Len() != 0 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}
	})

	t.Run("Test global warn", func(t *testing.T) {
		b := test.NewStringBuffer()
		l, err := New(
			WithLevel(LevelWarn),
			WithOutput(b),
			WithoutTimestamp(),
			AsDefault(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if DefaultLogger != l {
			t.Errorf("Expected default logger to be set")
		}

		Warn("a", "b")
		Warn("c", "d")

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
		b := test.NewStringBuffer()
		l, err := New(
			WithLevel(LevelError),
			WithOutput(b),
			WithoutTimestamp(),
			AsDefault(),
		)
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
		if DefaultLogger != l {
			t.Errorf("Expected default logger to be set")
		}

		Error("a", "b")
		Error("c", "d")

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
		{s: "", f: FormatLogfmt},
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
	t.Run("no logger in context", func(t *testing.T) {
		l := LoggerFromContext(context.Background())
		if l == nil {
			t.Errorf("Failed to create a logger")
		}
	})

	t.Run("insert, extract success", func(t *testing.T) {
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

		if b.Len() != 1 {
			t.Errorf("Unexpected buffer length: %d", b.Len())
		}

		expected := "level=info id=0000-111 bob=smith alice=jones\n"
		if !strings.Contains(b.Index(0), expected) {
			t.Errorf("\nExpected %sActual:  %s", expected, b.Index(0))
		}
	})
}
