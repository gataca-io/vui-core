package tools

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSliceTools_UniqueSlice(t *testing.T) {
	s := []string{"1", "2", "3", "2", "2", "4", "5", "0"}
	s1 := UniqueSlice(s)

	assert.Equal(t, len(s1), 6)
	assert.NotEqual(t, s, s1)
	assert.Contains(t, s1, "0")
	assert.Contains(t, s1, "1")
	assert.Contains(t, s1, "2")
	assert.Contains(t, s1, "3")
	assert.Contains(t, s1, "4")
	assert.Contains(t, s1, "5")
}

func TestSliceTools_Contains(t *testing.T) {
	s := []string{"1", "2", "3", "2", "2", "4", "5", "0"}
	s0 := Contains(s, "0")
	assert.Equal(t, s0, true)
	s1 := Contains(s, "1")
	assert.Equal(t, s1, true)
	s2 := Contains(s, "2")
	assert.Equal(t, s2, true)
	s3 := Contains(s, "3")
	assert.Equal(t, s3, true)
	s4 := Contains(s, "4")
	assert.Equal(t, s4, true)
	s5 := Contains(s, "5")
	assert.Equal(t, s5, true)
	s6 := Contains(s, "6")
	assert.Equal(t, s6, false)
}
