package maps_test

import (
	"sort"
	"testing"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/maps"
)

func TestKeys(t *testing.T) {
	var (
		empty   = map[string]string{}
		squares = map[int]int{1: 1, 2: 4, 3: 9, 4: 16, 5: 25}
		doubles = map[string]int{"1": 2, "2": 4, "3": 6, "4": 8, "5": 10}
	)

	assert.Compare(t, maps.Keys(empty), []string{})

	s1 := maps.Keys(squares)
	sort.Slice(s1, func(i, j int) bool { return s1[i] < s1[j] })
	assert.Compare(t, s1, []int{1, 2, 3, 4, 5})

	s2 := maps.Keys(doubles)
	sort.Slice(s2, func(i, j int) bool { return s2[i] < s2[j] })
	assert.Compare(t, s2, []string{"1", "2", "3", "4", "5"})
}

func TestMerge(t *testing.T) {
	var (
		existing = map[string]string{"en": "Banana"}
		updated  = map[string]string{"en": "Apple", "fr": "Pomme"}
	)

	m := maps.Merge(existing, updated)
	assert.Equal(t, len(m), 2)
	assert.Equal(t, m["en"], "Apple")
	assert.Equal(t, m["fr"], "Pomme")
}

func TestPretty(t *testing.T) {
	m := map[string]string{"a": "A A", "b": "B", "c": "C"}
	s := maps.Pretty(m)
	assert.Equal(t, s, "[a: A A, b: B, c: C]")
}
