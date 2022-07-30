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

package log

import (
	"testing"
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
		b    = NewStringBuffer()
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
