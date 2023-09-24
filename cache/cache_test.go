package cache

import (
	"testing"
	"time"

	"arcadium.dev/core/assert"
)

func TestCache(t *testing.T) {
	c := New[string, int](10 * time.Millisecond)

	assert.NotNil(t, c)
	assert.NotNil(t, c.cache)
	assert.Equal(t, len(c.cache), 0)
	assert.Equal(t, c.lifetime, 10*time.Millisecond)

	_, ok := c.Get("foobar")
	assert.False(t, ok)

	c.Set("zero", 0)
	c.Set("one", 1)
	c.Set("one", 11)
	c.Set("one", 1)

	assert.Equal(t, len(c.cache), 2)

	c.Set("two", 2)
	c.Set("three", 3)

	assert.Equal(t, len(c.cache), 4)

	keys := []string{"zero", "one", "two", "three"}
	i := 0
	for _, k := range keys {
		v, ok := c.Get(k)
		assert.True(t, ok)
		assert.Equal(t, i, v)
		i++
	}

	time.Sleep(10 * time.Millisecond)
	i = 0
	for _, k := range keys {
		_, ok := c.Get(k)
		assert.False(t, ok)
		i++

		c.Remove(k)
		assert.Equal(t, len(c.cache), len(keys)-i)
	}
	assert.Equal(t, len(c.cache), 0)
}
