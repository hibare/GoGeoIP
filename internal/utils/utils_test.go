package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStringInSlice(t *testing.T) {
	assert.True(t, SliceContains("1", []string{"1", "2", "3"}))
	assert.False(t, SliceContains("11", []string{"1", "2", "3"}))
	assert.True(t, SliceContains(2, []int{1, 2, 3}))
	assert.False(t, SliceContains(22, []int{1, 2, 3}))
	assert.True(t, SliceContains(3.3, []float64{1.1, 2.2, 3.3}))
	assert.False(t, SliceContains(33.3, []float64{1.1, 2.2, 3.3}))
}
