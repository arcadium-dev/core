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

// Package slices provides useful functions with slices of any type.
package slices // import "arcadium.dev/core/slices"

import (
	"fmt"
	"sort"
	"strings"

	"golang.org/x/exp/constraints"
)

// Contains returns true if the given slices s contains k.
func Contains[S ~[]E, E comparable](s S, e E) bool {
	for _, c := range s {
		if c == e {
			return true
		}
	}
	return false
}

// Intersection returns a slice of elements that are both in s1 and s2.
// Order matters, the elements from s1 are compared to s2, so the resultant
// slice will be a subset of s1. If s1 has duplicates entries that
// intersect, they will be in the result.
func Intersection[S ~[]E, E comparable](s1, s2 S) S {
	if len(s1) == 0 || len(s2) == 0 {
		return make(S, 0)
	}

	result := make(S, 0)
	for _, e := range s1 {
		if Contains(s2, e) {
			result = append(result, e)
		}
	}
	return result
}

// Diff returns a slice of elements are in s1 but not in s2.
func Diff[S ~[]E, E comparable](s1, s2 S) S {
	if len(s1) == 0 {
		return make(S, 0)
	}
	if len(s2) == 0 {
		return s1
	}

	result := make(S, 0)
	for _, e := range s1 {
		if !Contains(s2, e) {
			result = append(result, e)
		}
	}
	return result
}

// Equal returns true if the slices are the same.
func Equal[S ~[]E, E comparable](l, r S) bool {
	ln := len(l)
	if ln != len(r) {
		return false
	}
	for i := 0; i < ln; i++ {
		if l[i] != r[i] {
			return false
		}
	}
	return true
}

// Dedup returns a sorted slice with duplicate entries removed.
func Dedup[S ~[]E, E constraints.Ordered](s S) S {
	result := make(S, 0)
	sort.Slice(s, func(i, j int) bool { return s[i] < s[j] })
	j := 0
	for i, e := range s {
		if i == 0 || e != result[j-1] {
			result = append(result, e)
			j++
		}
	}
	return result
}

// Last returns the last element of the slice.
func Last[S ~[]E, E any](s S) E {
	var empty E
	l := len(s)
	if l == 0 {
		return empty
	}
	return s[l-1]
}

// Merge returns the result of the merged slices.
func Merge[S ~[]E, E comparable](l, r S) S {
	result := make(S, len(l))
	copy(result, l)

	for _, e := range r {
		if !Contains(result, e) {
			result = append(result, e)
		}
	}

	return result
}

// Pretty prints the slice.
func Pretty(s []string) string {
	p := fmt.Sprintf("%q", s)
	p = strings.ReplaceAll(p, `" "`, `, `)
	p = strings.ReplaceAll(p, `["`, `[`)
	p = strings.ReplaceAll(p, `"]`, `]`)
	return p
}
