package collections

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndexOfWithEmptySlice(t *testing.T) {
	slice := []string{}
	index := IndexOf(slice, "test")

	assert.Equal(t, -1, index)
}

func TestIndexOfSingleItemSliceMatch(t *testing.T) {
	slice := []string{"test"}
	index := IndexOf(slice, "test")

	assert.Equal(t, 0, index)
}

func TestIndexOfMultipleItemSliceMatch(t *testing.T) {
	slice := []string{"test", "123", "mic check", "1, 2"}
	index := IndexOf(slice, "mic check")

	assert.Equal(t, 2, index)
}