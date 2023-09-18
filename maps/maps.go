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

// Package maps provides useful functions with maps of any type.
package maps // import "arcadium.dev/core/maps"

import (
	"fmt"
	"strings"
)

// Keys returns the keys of the map m. The keys will be in an indeterminate order.
func Keys[M ~map[K]V, K comparable, V any](m M) []K {
	result := make([]K, len(m))
	i := 0
	for k := range m {
		result[i] = k
		i++
	}
	return result
}

// Merge returns the result of the merged maps. If a key/value pair exists in both maps, l and r, the
// key/value pair from r will take precedent.
func Merge[M ~map[K]V, K comparable, V any](l, r M) M {
	result := make(M)

	for k, v := range r {
		result[k] = v
	}
	for k, v := range l {
		if _, ok := result[k]; !ok {
			result[k] = v
		}
	}
	return result
}

// Pretty prints the map.
func Pretty(m map[string]string) string {
	p := fmt.Sprintf("%q", m)
	p = strings.ReplaceAll(p, `" "`, `, `)
	p = strings.ReplaceAll(p, `":"`, `: `)
	p = strings.ReplaceAll(p, `map["`, `[`)
	p = strings.ReplaceAll(p, `"]`, `]`)
	return p
}
