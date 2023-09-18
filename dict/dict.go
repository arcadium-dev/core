// Copyright 2023 arcadium.dev <info@arcadium.dev>
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

// Package dict provides a data structure that wraps a map and allows
// repeatable traversals. The ordering is based on the order which the
// key / value pairs is set.
package dict // import "arcadium.dev/core/dict"

// Dict wraps a map and allows for repeatable traversals, via the Keys method.
type Dict[K comparable, V any] struct {
	m map[K]V
	k []K
}

// New creates a new dict.
func New[K comparable, V any]() *Dict[K, V] {
	return &Dict[K, V]{
		m: map[K]V{},
		k: make([]K, 0),
	}
}

// Set creates a new entry in the Dict for the key k, and the value v.  The
// ordering of the Dict is based on the order in which the key / value is Set.
// If the key / value pair already exists, it will be replaced and the original
// ordering is preserved.
func (d *Dict[K, V]) Set(k K, v V) {
	if _, ok := d.m[k]; !ok {
		d.k = append(d.k, k)
	}
	d.m[k] = v
}

// Len returns the length of the Dict.
func (d *Dict[K, V]) Len() int {
	return len(d.k)
}

// Keys returns the keys of the Dict. The order of the slice is the order in
// which the values were set.
func (d Dict[K, V]) Keys() []K {
	keys := make([]K, len(d.k))
	copy(keys, d.k)
	return keys
}

// Key returns the key at index i.
func (d *Dict[K, V]) Key(i int) K {
	return d.k[i]
}

// Val returns the value for key k.
func (d *Dict[K, V]) Val(k K) V {
	return d.m[k]
}
