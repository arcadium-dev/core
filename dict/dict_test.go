package dict

import (
	"testing"

	"arcadium.dev/core/assert"
)

func TestDict(t *testing.T) {
	d := New[string, int]()

	assert.NotNil(t, d)
	assert.NotNil(t, d.m)
	assert.NotNil(t, d.k)
	assert.Equal(t, len(d.m), 0)
	assert.Equal(t, len(d.k), 0)

	d.Set("zero", 0)

	d.Set("one", 1)
	d.Set("one", 11)
	d.Set("one", 1)

	assert.Equal(t, d.Len(), 2)

	d.Set("two", 2)
	d.Set("three", 3)

	assert.Equal(t, len(d.m), 4)
	assert.Equal(t, len(d.k), 4)
	assert.Equal(t, d.Len(), 4)

	keys := d.Keys()
	assert.Equal(t, len(keys), 4)

	d.Set("four", 4)
	assert.Equal(t, len(keys), 4)
	assert.Equal(t, d.Len(), 5)

	i := 0
	for _, k := range keys {
		assert.Equal(t, k, d.Key(i))
		assert.Equal(t, d.Val(k), i)
		i++
	}
}
