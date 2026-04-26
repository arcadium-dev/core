// Copyright 2021-2026 arcadium.dev <info@arcadium.dev>
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

// Package require provides a set of tools for use with the Go testing package.
package require // import "arcadium.dev/core/require"

import (
	"reflect"
	"testing"
)

// Equal requires that the expected value is equal to the actual value.
// Actual and expected must be comparable.
func Equal[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual != expected {
		t.Fatalf("\nExpected: %+v\nActual:   %v", expected, actual)
	}
}

// NotEqual requires that the expected value is not equal to the actual value.
// Actual and expected must be comparable.
func NotEqual[T comparable](t *testing.T, actual, expected T) {
	t.Helper()
	if actual == expected {
		t.Fatalf("Expected different values: %+v", actual)
	}
}

// Nil requires that the value of the given object is nil.
func Nil(t *testing.T, object any) {
	t.Helper()

	if object == nil {
		return
	}

	value := reflect.ValueOf(object)
	switch value.Kind() {
	case
		reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice, reflect.UnsafePointer:

		if value.IsNil() {
			return
		}
	}

	t.Fatalf("Unexpected non-nil value: %+v", object)
}

// NotNil requires that the value of the given object is not nil.
func NotNil(t *testing.T, object any) {
	t.Helper()
	if object != nil {
		return
	}

	value := reflect.ValueOf(object)
	switch value.Kind() {
	case
		reflect.Chan, reflect.Func,
		reflect.Interface, reflect.Map,
		reflect.Pointer, reflect.Slice, reflect.UnsafePointer:

		if !value.IsNil() {
			return
		}
	}

	t.Fatalf("Unexpected nil value: %+v", object)
}

// True requires that the given value is true.
func True(t *testing.T, value bool) {
	t.Helper()
	if !value {
		t.Fatal("Expected value to be true")
	}
}
