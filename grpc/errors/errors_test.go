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

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type testError string

func (e testError) Error() string { return string(e) }

func TestError(t *testing.T) {
	t.Parallel()

	t.Run("Test New", func(t *testing.T) {
		t.Parallel()

		err := New(codes.Unknown, "test error")
		if err.Error() != "test error" {
			t.Errorf("Unexpected error: %s", err)
		}
		code := status.Code(err)
		if code != codes.Unknown {
			t.Errorf("Unexpected code: %s", code)
		}
	})

	t.Run("Test WithCode", func(t *testing.T) {
		t.Parallel()

		err := WithCode(testError("test error"), codes.Unimplemented)
		if err.Error() != "test error" {
			t.Errorf("Unexpected error: %s", err)
		}
		code := status.Code(err)
		if code != codes.Unimplemented {
			t.Errorf("Unexpected code: %s", code)
		}
	})

	t.Run("Test Errorf", func(t *testing.T) {
		t.Parallel()

		err := Errorf(codes.NotFound, "%s %s", "test", "error")
		if err.Error() != "test error" {
			t.Errorf("Unexpected err: %s", err)
		}
		code := status.Code(err)
		if code != codes.NotFound {
			t.Errorf("Unexpected code: %s", code)
		}
	})

	t.Run("Test Wrap", func(t *testing.T) {
		t.Parallel()

		err := Wrap(nil, "test error")
		if err != nil {
			t.Errorf("Unexpected err: %s", err)
		}
	})

	t.Run("Test Wrapf", func(t *testing.T) {
		t.Parallel()

		err := Wrapf(InvalidArgument, "test error")
		if err.Error() != "InvalidArgument: test error" {
			t.Errorf("Unexpected err: %s", err)
		}
		code := status.Code(err)
		if code != codes.InvalidArgument {
			t.Errorf("Unexpected code: %s", code)
		}
	})

	t.Run("Test Unwrap", func(t *testing.T) {
		t.Parallel()

		err := Unwrap(Wrap(AlreadyExists, "wrapped"))
		if err.Error() != "AlreadyExists" {
			t.Errorf("Unexpected err: %s", err)
		}
		code := status.Code(err)
		if code != codes.AlreadyExists {
			t.Errorf("Unexpected code: %s", code)
		}
	})

	t.Run("Test Is", func(t *testing.T) {
		t.Parallel()

		err := Wrap(Internal, "wrapped")
		if !Is(err, Internal) {
			t.Errorf("Unexpected err: %s", err)
		}
	})

	t.Run("Test As", func(t *testing.T) {
		t.Parallel()

		var e testError
		err := Wrapf(testError("foo bar"), "wrapped")
		if !As(err, &e) {
			t.Errorf("Unexpected err: %s %s", err, e)
		}
		if e.Error() != "foo bar" {
			t.Errorf("Unexpected err: %s", e)
		}
	})

	t.Run("Test GRPCStatus", func(t *testing.T) {
		t.Parallel()

		_, ok := status.FromError(InvalidArgument)
		if !ok {
			t.Error("failed type assertion")
		}

		code := status.Code(InvalidArgument)
		if code != codes.InvalidArgument {
			t.Errorf("unexpected code: %s", code)
		}
	})
}

func TestSentinelErrors(t *testing.T) {
	t.Parallel()

	var table = []struct {
		err  error
		code codes.Code
	}{
		{Unknown, codes.Unknown},
		{InvalidArgument, codes.InvalidArgument},
		{NotFound, codes.NotFound},
		{AlreadyExists, codes.AlreadyExists},
		{Unimplemented, codes.Unimplemented},
		{Internal, codes.Internal},
	}

	for _, entry := range table {
		_, ok := status.FromError(entry.err)
		if !ok {
			t.Error("failed type assertion")
		}

		code := status.Code(entry.err)
		if code != entry.code {
			t.Errorf("unexpected code: %s", code)
		}

		msg := entry.err.Error()
		if msg != entry.code.String() {
			t.Errorf("unsexpected msg: %s", msg)
		}
	}
}
