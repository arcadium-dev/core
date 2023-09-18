package slices_test

import (
	"testing"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/slices"
)

var (
	emptyStrings = []string{}
	emptyInts    = []int{}

	alphabet = []string{
		"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m",
		"n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z",
	}
	digits = []int{
		0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
	}

	foobar = []string{
		"f", "o", "o", "b", "a", "r", "B", "A", "z",
	}
	teel = []int{7, 3, 42, 3, 1, 73}
)

func TestContains(t *testing.T) {
	assert.False(t, slices.Contains(emptyStrings, "foo"))
	assert.False(t, slices.Contains(emptyStrings, ""))
	assert.False(t, slices.Contains(emptyInts, -1))

	assert.False(t, slices.Contains(alphabet, "0"))
	assert.False(t, slices.Contains(digits, 42))

	assert.True(t, slices.Contains(alphabet, "a"))
	assert.True(t, slices.Contains(alphabet, "r"))
	assert.True(t, slices.Contains(alphabet, "z"))
	assert.True(t, slices.Contains(digits, 0))
	assert.True(t, slices.Contains(digits, 4))
	assert.True(t, slices.Contains(digits, 9))
}

func TestIntersection(t *testing.T) {
	assert.Compare(t, slices.Intersection(alphabet, emptyStrings), emptyStrings)
	assert.Compare(t, slices.Intersection(emptyStrings, alphabet), emptyStrings)

	assert.Compare(t, slices.Intersection(emptyInts, digits), emptyInts)
	assert.Compare(t, slices.Intersection(digits, emptyInts), emptyInts)

	assert.Compare(t, slices.Intersection(alphabet, foobar), []string{"a", "b", "f", "o", "r", "z"})
	assert.Compare(t, slices.Intersection(foobar, alphabet), []string{"f", "o", "o", "b", "a", "r", "z"})

	assert.Compare(t, slices.Intersection(digits, teel), []int{1, 3, 7})
	assert.Compare(t, slices.Intersection(teel, digits), []int{7, 3, 3, 1})
}

func TestDiff(t *testing.T) {
	assert.Compare(t, slices.Diff(alphabet, emptyStrings), alphabet)
	assert.Compare(t, slices.Diff(emptyStrings, alphabet), emptyStrings)

	assert.Compare(t, slices.Diff(digits, emptyInts), digits)
	assert.Compare(t, slices.Diff(emptyInts, digits), emptyInts)

	assert.Compare(t, slices.Diff(alphabet, foobar), []string{
		"c", "d", "e", "g", "h", "i", "j", "k", "l", "m", "n", "p", "q", "s", "t", "u", "v", "w", "x", "y",
	})
	assert.Compare(t, slices.Diff(foobar, alphabet), []string{"B", "A"})

	assert.Compare(t, slices.Diff(digits, teel), []int{0, 2, 4, 5, 6, 8, 9})
	assert.Compare(t, slices.Diff(teel, digits), []int{42, 73})
}

func TestEqual(t *testing.T) {
	var (
		a = []string{"a", "b", "c"}
		b = []string{"c", "b", "a"}
		c = []string{"a", "b", "c"}
	)

	assert.False(t, slices.Equal(foobar, alphabet))
	assert.False(t, slices.Equal(a, b))
	assert.True(t, slices.Equal(a, a))
	assert.True(t, slices.Equal(a, c))

	var (
		d = []int{1, 2, 3}
		e = []int{3, 2, 1}
		f = []int{1, 2, 3}
	)
	assert.False(t, slices.Equal(digits, teel))
	assert.False(t, slices.Equal(d, e))
	assert.True(t, slices.Equal(d, d))
	assert.True(t, slices.Equal(d, f))
}

func TestDedup(t *testing.T) {
	assert.Compare(t, slices.Dedup(emptyInts), emptyInts)
	assert.Compare(t, slices.Dedup(foobar), []string{"A", "B", "a", "b", "f", "o", "r", "z"})
	assert.Compare(t, slices.Dedup(alphabet), alphabet)
	assert.Compare(t, slices.Dedup(digits), digits)
	assert.Compare(t, slices.Dedup(teel), []int{1, 3, 7, 42, 73})
}

func TestLast(t *testing.T) {
	assert.Equal(t, 0, slices.Last(emptyInts))
	assert.Equal(t, "z", slices.Last([]string{"A", "B", "a", "b", "f", "o", "r", "z"}))
	assert.Equal(t, "z", slices.Last(alphabet))
	assert.Equal(t, 9, slices.Last(digits))
	assert.Equal(t, 73, slices.Last(teel))
}

func TestMerge(t *testing.T) {
	var (
		empty = []string{}
		a     = []string{"a", "aa", "A", "AA"}
		b     = []string{"b", "bb", "B", "BB"}
		ab    = []string{"a", "b", "aa", "bb", "AA", "BB", "AB"}
	)

	t.Run("merge empty", func(t *testing.T) {
		s := slices.Merge(empty, empty)
		assert.NotNil(t, s)
		assert.Equal(t, len(s), 0)
	})

	t.Run("no overlap", func(t *testing.T) {
		s1 := slices.Merge(a, b)
		assert.Equal(t, len(s1), len(a)+len(b))
		assert.Compare(t, s1, append(a, b...))

		s2 := slices.Merge(b, a)
		assert.Equal(t, len(s2), len(a)+len(b))
		assert.Compare(t, s2, append(b, a...))
	})

	t.Run("full overlap", func(t *testing.T) {
		s := slices.Merge(a, a)
		assert.Compare(t, s, a)
	})

	t.Run("partial overlap", func(t *testing.T) {
		s1 := slices.Merge(a, ab)
		assert.Equal(t, len(s1), len(a)+4)
		assert.Compare(t, s1, append(a, []string{"b", "bb", "BB", "AB"}...))

		s2 := slices.Merge(b, ab)
		assert.Equal(t, len(s2), len(b)+4)
		assert.Compare(t, s2, append(b, []string{"a", "aa", "AA", "AB"}...))
	})
}

func TestPretty(t *testing.T) {
	s := []string{"this is a test", "a", "b"}
	assert.Equal(t, "[this is a test, a, b]", slices.Pretty(s))
}
