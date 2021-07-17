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

func TestGameError(t *testing.T) {
	// Test Errorf()
	err := Errorf(codes.Unimplemented, "%s%s", "foo", "bar")
	if err.Error() != "foobar" {
		t.Errorf("unexpected error: %+v", err)
	}

	code := status.Code(err)
	if code != codes.Unimplemented {
		t.Errorf("unexpected code: %s", code)
	}

	// Test wrapping a nil error
	err = Wrapf(nil, "%s", "Player ID")

	if err != nil {
		t.Errorf("unexpected wrapped nil error: %s", err)
	}

	// Test wrapping a valid error
	err = Wrapf(InvalidArgument, "%s", "Player ID")

	// Test Error()
	msg := err.Error()
	if msg != "InvalidArgument: Player ID" {
		t.Errorf("unexpected error message: %s", msg)
	}

	// Test Unwrap()
	base := Unwrap(err)
	if base != InvalidArgument {
		t.Errorf("unexpected unwrapped error: %s", base)
	}

	// Test Is()
	if !Is(err, InvalidArgument) {
		t.Error("should be an invalid player id error")
	}

	// Test As()
	var e *gerror
	if !As(err, &e) {
		t.Error("should be a game error")
	}
	if e != err {
		t.Error("should be an invalid player id error")
	}

	// Test GRPCStatus()
	_, ok := status.FromError(err)
	if !ok {
		t.Error("failed type assertion")
	}

	code = status.Code(err)
	if code != codes.InvalidArgument {
		t.Errorf("unexpected code: %s", code)
	}
}

func TestGameBaseErrors(t *testing.T) {
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
