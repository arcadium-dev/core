// Copyright 2022 arcadium.dev <info@arcadium.dev>
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

package http

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	cerrors "arcadium.dev/core/errors"
)

func TestResponse(t *testing.T) {
	ctx := context.Background()

	t.Run("nil error", func(t *testing.T) {
		Response(nil, nil, nil)
	})

	t.Run("invalid argument error", func(t *testing.T) {
		w := httptest.NewRecorder()
		Response(ctx, w, cerrors.ErrInvalidArgument)

		checkHeader(t, w, "Content-Type", "application/json")
		checkHeader(t, w, "X-Content-Type-Options", "nosniff")

		if w.Code != http.StatusBadRequest {
			t.Errorf("Unexpected status: %d", w.Code)
		}
		expected := `{"error":{"status":400,"detail":"invalid argument"}}` + "\n"
		if w.Body.String() != expected {
			t.Errorf("\nExpected body %s\nActual body   %s", expected, w.Body.String())
		}
	})

	t.Run("not found error", func(t *testing.T) {
		w := httptest.NewRecorder()
		Response(ctx, w, cerrors.ErrNotFound)

		checkHeader(t, w, "Content-Type", "application/json")
		checkHeader(t, w, "X-Content-Type-Options", "nosniff")

		if w.Code != http.StatusNotFound {
			t.Errorf("Unexpected status: %d", w.Code)
		}
		expected := `{"error":{"status":404,"detail":"not found"}}` + "\n"
		if w.Body.String() != expected {
			t.Errorf("\nExpected body %s\nActual body   %s", expected, w.Body.String())
		}
	})

	t.Run("already exists error", func(t *testing.T) {
		w := httptest.NewRecorder()
		Response(ctx, w, cerrors.ErrAlreadyExists)

		checkHeader(t, w, "Content-Type", "application/json")
		checkHeader(t, w, "X-Content-Type-Options", "nosniff")

		if w.Code != http.StatusConflict {
			t.Errorf("Unexpected status: %d", w.Code)
		}
		expected := `{"error":{"status":409,"detail":"already exists"}}` + "\n"
		if w.Body.String() != expected {
			t.Errorf("\nExpected body %s\nActual body   %s", expected, w.Body.String())
		}
	})

	t.Run("not implemented error", func(t *testing.T) {
		w := httptest.NewRecorder()
		Response(ctx, w, cerrors.ErrNotImplemented)

		checkHeader(t, w, "Content-Type", "application/json")
		checkHeader(t, w, "X-Content-Type-Options", "nosniff")

		if w.Code != http.StatusNotImplemented {
			t.Errorf("Unexpected status: %d", w.Code)
		}
		expected := `{"error":{"status":501,"detail":"not implemented"}}` + "\n"
		if w.Body.String() != expected {
			t.Errorf("\nExpected body %s\nActual body   %s", expected, w.Body.String())
		}
	})

	t.Run("internal server error", func(t *testing.T) {
		w := httptest.NewRecorder()
		Response(ctx, w, cerrors.ErrInternal)

		checkHeader(t, w, "Content-Type", "application/json")
		checkHeader(t, w, "X-Content-Type-Options", "nosniff")

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Unexpected status: %d", w.Code)
		}
		expected := `{"error":{"status":500,"detail":"internal error"}}` + "\n"
		if w.Body.String() != expected {
			t.Errorf("\nExpected body %s\nActual body   %s", expected, w.Body.String())
		}
	})

	t.Run("unknown error", func(t *testing.T) {
		w := httptest.NewRecorder()
		Response(ctx, w, errors.New("unknown error"))

		checkHeader(t, w, "Content-Type", "application/json")
		checkHeader(t, w, "X-Content-Type-Options", "nosniff")

		if w.Code != http.StatusInternalServerError {
			t.Errorf("Unexpected status: %d", w.Code)
		}
		expected := `{"error":{"status":500,"detail":"unknown error"}}` + "\n"
		if w.Body.String() != expected {
			t.Errorf("\nExpected body %s\nActual body   %s", expected, w.Body.String())
		}
	})
}

func TestResponseError(t *testing.T) {
	expected := `status=200, detail="foobar"`
	err := responseError{Status: http.StatusOK, Detail: "foobar"}
	if err.Error() != expected {
		t.Errorf("\nExpected error %s\nActual error   %s", expected, err)
	}
}

func checkHeader(t *testing.T, w http.ResponseWriter, header, expected string) {
	if w.Header().Get(header) != expected {
		t.Errorf("Unexpected %s header: %s", header, w.Header().Get(header))
	}
}
