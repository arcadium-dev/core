// Copyright 2021-2023 arcadium.dev <info@arcadium.dev>
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

// Package assert provides a set of tools for use with the Go testing package.
package assert // import "arcadium.dev/core/assert"

import (
	"encoding/json"
	"errors"
	"io"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/go-cmp/cmp"

	"arcadium.dev/core/http/server"
)

// Contains asserts that the expected string is contained in the actual string.
func Contains(t *testing.T, actual, expected string) {
	t.Helper()
	if !strings.Contains(actual, expected) {
		t.Errorf("\nActual:   %s\nExpected: %s", actual, expected)
	}
}

// Equal asserts that the expected value is equal to the actual value.
// Actual and expected must be comparable. To compare non-comparable
// types use Compare.
func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Errorf("\nActual:   %+v\nExpected: %+v", actual, expected)
	}
}

// NotEqual asserts that the expected value is not equal to the actual value.
// Actual and expected must be comparable. To compare non-comparable types use
// Compare.
func NotEqual[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual == expected {
		t.Errorf("Expected different values: %+v", actual)
	}
}

// Compare asserts that the expected value is equal to the actual value.
// Actual and expected need not be comparable.
func Compare[T any](t *testing.T, actual, expected T, opts ...cmp.Option) {
	t.Helper()
	if !cmp.Equal(actual, expected, opts...) {
		t.Errorf("\nActual:   %+v\nExpected: %+v", actual, expected)
	}
}

// Error asserts that the error string from err matches the expected.
func Error(t *testing.T, err error, expected string) {
	t.Helper()
	if err == nil {
		t.Fatal("Expected an error")
	}
	if expected != err.Error() {
		t.Errorf("\nActual error:   %s\nExpected error: %s", err, expected)
	}
}

// IsError asserts that the actual error matches the expected error.
func IsError(t *testing.T, actual, expected error) {
	t.Helper()
	if actual == nil {
		t.Fatal("Expected an error")
	}
	if !errors.Is(actual, expected) {
		t.Errorf("\nActual:   %s\nExpected: %s", actual, expected)
	}
}

func ResponseError(t *testing.T, w *httptest.ResponseRecorder, status int, errMsg string) {
	resp := w.Result()
	Equal(t, resp.StatusCode, status)

	body, err := io.ReadAll(resp.Body)
	Nil(t, err)
	defer resp.Body.Close()

	var respErr server.ResponseError
	err = json.Unmarshal(body, &respErr)
	Nil(t, err)

	Contains(t, respErr.Detail, errMsg)
	Equal(t, respErr.Status, status)
}

// MockExpectationsMet asserts that the expectations for the given mock were
// met.
func MockExpectationsMet(t *testing.T, mock sqlmock.Sqlmock) {
	t.Helper()
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Unexpected error: %s", err)
	}
}

// Nil asserts that the value of the given object is nil.
func Nil(t *testing.T, object any) {
	t.Helper()

	value := reflect.ValueOf(object)
	kind := value.Kind()

	for _, k := range []reflect.Kind{reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice} {
		if k == kind {
			if value.IsNil() {
				return
			}
		}
	}
	if object == nil {
		return
	}

	t.Errorf("Unexpected non-nil value: %+v", object)
}

// NotNil asserts that the value of the given object is not nil.
func NotNil(t *testing.T, object any) {
	t.Helper()

	value := reflect.ValueOf(object)
	kind := value.Kind()

	for _, k := range []reflect.Kind{reflect.Chan, reflect.Func, reflect.Interface, reflect.Map, reflect.Ptr, reflect.Slice} {
		if k == kind {
			if !value.IsNil() {
				return
			}
		}
	}
	if object != nil {
		return
	}

	t.Errorf("Unexpected nil value: %+v", object)
}

// True asserts that the given value is true.
func True(t *testing.T, value bool) {
	t.Helper()
	if !value {
		t.Error("Expected value to be true")
	}
}

// False asserts that the given value is false.
func False(t *testing.T, value bool) {
	t.Helper()
	if value {
		t.Error("Expected value to be false")
	}
}
