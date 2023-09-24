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

// Package cache provides a simple cache.
package cache // import "arcadium.dev/core/cache"

import (
	"log"
	"sync"
	"time"
)

// Cache is a thread safe key/value map where the entries can expire.
type Cache[K comparable, V any] struct {
	mu       sync.RWMutex
	cache    map[K]*entry[V]
	lifetime time.Duration
}

// New creates a new cache where the lifetime is the duration of a
// valid cache entry, after which the entry is considered expired.
func New[K comparable, V any](lifetime time.Duration) *Cache[K, V] {
	if lifetime <= 0 {
		log.Fatalf("Failed to create cache, invalid lifefime: %d", lifetime)
	}
	return &Cache[K, V]{
		cache:    make(map[K]*entry[V]),
		lifetime: lifetime,
	}
}

// Get returns a copy of of the value V associated with the key, and the status
// of that value. A status of true indicates the value can be used. A status of
// false indicates the value is unusable, being either expired or not in the
// cache.
func (c *Cache[K, V]) Get(key K) (V, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var emptyValue V
	e, ok := c.cache[key]
	if !ok {
		return emptyValue, false
	}
	return e.get()
}

// Set will create or update the value associated for the given key.
func (c *Cache[K, V]) Set(key K, value V) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cache[key] = &entry[V]{value: value, expiry: time.Now().UTC().Add(c.lifetime)}
}

// Remove deletes value from the cache. If the value isn't present in the
// cache, this is a no-op.
func (c *Cache[K, V]) Remove(key K) {
	c.mu.Lock()
	defer c.mu.Unlock()
	delete(c.cache, key)
}

type entry[V any] struct {
	value  V
	expiry time.Time
}

func (e entry[V]) get() (V, bool) {
	return e.value, e.expiry.After(time.Now().UTC())
}
