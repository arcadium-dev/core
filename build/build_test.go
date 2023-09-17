package build_test

import (
	"testing"

	"arcadium.dev/core/assert"
	"arcadium.dev/core/build"
)

func TestString(t *testing.T) {
	info := build.Info("Testing", "Version", "Branch", "Commit", "Date")
	info.Go = "Go"
	assert.Equal(t, info.String(), "Testing Version (branch: Branch, commit: Commit, date: Date, go: Go)")
}
