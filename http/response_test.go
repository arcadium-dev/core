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
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestResponse(t *testing.T) {
	t.Run("no errors", func(t *testing.T) {
		w := httptest.NewRecorder()
		Response(w)

		checkHeader(t, w, "Content-Type", "application/json")
		checkHeader(t, w, "X-Content-Type-Options", "nosniff")
		if w.Code != http.StatusOK {
			t.Errorf("Unexpected status: %d", w.Code)
		}
		expected := ""
		if w.Body.String() != expected {
			t.Errorf("\nExpected body %s\nActual body   %s", w.Body.String(), expected)
		}
	})

	t.Run("single error", func(t *testing.T) {
		w := httptest.NewRecorder()
		Response(w, BadRequestError(errors.New("invalid argument")))

		checkHeader(t, w, "Content-Type", "application/json")
		checkHeader(t, w, "X-Content-Type-Options", "nosniff")
		if w.Code != http.StatusBadRequest {
			t.Errorf("Unexpected status: %d", w.Code)
		}
		expected := `{"errors":[{"status":400,"detail":"invalid argument"}]}` + "\n"
		if w.Body.String() != expected {
			t.Errorf("\nExpected body %s\nActual body   %s", expected, w.Body.String())
		}
	})

	t.Run("mutliple errors", func(t *testing.T) {
		w := httptest.NewRecorder()
		Response(w,
			NotFoundError(errors.New("not found")),
			ConflictError(errors.New("already exists")),
			InternalServerError(errors.New("internal error")),
			NotImplementedError(errors.New("not implemented")),
		)

		checkHeader(t, w, "Content-Type", "application/json")
		checkHeader(t, w, "X-Content-Type-Options", "nosniff")
		if w.Code != http.StatusNotFound {
			t.Errorf("Unexpected status: %d", w.Code)
		}
		expected := `{"errors":[{"status":404,"detail":"not found"},{"status":409,"detail":"already exists"},{"status":500,"detail":"internal error"},{"status":501,"detail":"not implemented"}]}` + "\n"
		if w.Body.String() != expected {
			t.Errorf("\nExpected body %s\nActual body   %s", expected, w.Body.String())
		}
	})

	t.Run("invalid status", func(t *testing.T) {
		w := httptest.NewRecorder()
		Response(w, ResponseError{Status: 1024, Detail: "foobar"})

		checkHeader(t, w, "Content-Type", "application/json")
		checkHeader(t, w, "X-Content-Type-Options", "nosniff")
		if w.Code != http.StatusInternalServerError {
			t.Errorf("Unexpected status: %d", w.Code)
		}
		expected := `{"errors":[{"status":500,"detail":"foobar"}]}` + "\n"
		if w.Body.String() != expected {
			t.Errorf("\nExpected body %s\nActual body   %s", expected, w.Body.String())
		}
	})
}

func TestResponseError(t *testing.T) {
	expected := `status=200, detail="foobar"`
	err := ResponseError{Status: http.StatusOK, Detail: "foobar"}
	if err.Error() != expected {
		t.Errorf("\nExpected error %s\nActual error   %s", expected, err)
	}
}

func checkHeader(t *testing.T, w http.ResponseWriter, header, expected string) {
	if w.Header().Get(header) != expected {
		t.Errorf("Unexpected %s header: %s", header, w.Header().Get(header))
	}
}
