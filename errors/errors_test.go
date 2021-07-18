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

package errors

import (
	"testing"
)

type testError string

func (e testError) Error() string { return string(e) }

var (
	Err        = New("error")
	ErrWrapped = Wrap(Err, "wrapped")

	ErrTest        = testError("test error")
	ErrTestWrapped = Wrap(ErrTest, "wrapped")
)

func TestError(t *testing.T) {
	t.Parallel()

	t.Run("Test New", func(t *testing.T) {
		t.Parallel()

		err := New("test error")
		if err.Error() != "test error" {
			t.Errorf("Unexpected error: %+v", err)
		}
	})

	t.Run("Test Errorf", func(t *testing.T) {
		t.Parallel()

		err := Errorf("%s %s", "test", "error")
		if err.Error() != "test error" {
			t.Errorf("Unexpected error: %+v", err)
		}
	})

	t.Run("Test Wrap", func(t *testing.T) {
		t.Parallel()

		err := Wrapf(nil, "%s", "test error")
		if err != nil {
			t.Errorf("Unexpected error: %s", err)
		}
	})

	t.Run("Test Wrapf", func(t *testing.T) {
		t.Parallel()

		err := Wrapf(New("first"), "%s", "second")
		if err.Error() != "first: second" {
			t.Errorf("Unexpected error: %s", err)
		}
	})

	t.Run("Test Unwrap", func(t *testing.T) {
		t.Parallel()

		err := Unwrap(ErrWrapped)
		if err == nil {
			t.Errorf("Expected an error")
		}
		if err != Err {
			t.Errorf("Unexpected error: %s", err)
		}
	})

	t.Run("Test Is", func(t *testing.T) {
		t.Parallel()

		if !Is(ErrWrapped, Err) {
			t.Errorf("Expected Err")
		}
	})

	t.Run("Test As", func(t *testing.T) {
		t.Parallel()

		var err testError
		if !As(ErrTestWrapped, &err) || err != ErrTest {
			t.Errorf("Unexpected errror: %s", err)
		}
	})
}
